package pinata_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stafiprotocol/eth-lsd-relay/pkg/destorage/pinata"
	"github.com/stretchr/testify/assert"
)

func TestUploadAndDownload(t *testing.T) {
	pinataApikey := os.Getenv("PINATA_APIKEY")
	pinataEndpoint := os.Getenv("PINATA_ENDPOINT")
	client, err := pinata.NewClient(pinata.Config{Endpoint: pinataEndpoint, Apikey: pinataApikey})
	assert.Nil(t, err)

	fileName := "hello-world.txt"
	fileContent := []byte("Hello World! - " + time.Now().String())
	cid, err := client.UploadFile(fileContent, fileName)
	assert.Nil(t, err)
	assert.NotEmpty(t, cid)
	fmt.Println("file name:", fileName)
	fmt.Println("cid:", cid)

	downloadContent, err := client.DownloadFile(cid, fileName)
	assert.Nil(t, err)
	assert.Equal(t, string(fileContent), string(downloadContent))
	fmt.Println("content:", string(downloadContent))
}

func TestUnpinFilesCreatedBefore(t *testing.T) {
	pinataApikey := os.Getenv("PINATA_APIKEY")
	pinataEndpoint := os.Getenv("PINATA_ENDPOINT")
	client, err := pinata.NewClient(pinata.Config{Endpoint: pinataEndpoint, Apikey: pinataApikey})
	assert.Nil(t, err)

	count, err := client.UnpinFilesCreatedBefore(time.Now().AddDate(0, -180, 0))
	assert.Nil(t, err)
	fmt.Println("deleted count:", count)
}

func TestDefaultGatewayFormat(t *testing.T) {
	cid := "bafybeiajiqzfn2o4qnlodg5yfuqtqzjxivo6vequroxwnpcmdokazmhs74"
	fileName := "0x79bb3a0ee435f957ce4f54ee8c3cfadc7278da0c-rewards-369-331725.json"

	url := fmt.Sprintf(pinata.DefaultGateway, cid, fileName)
	assert.Equal(t,
		"https://gateway.pinata.cloud/ipfs/bafybeiajiqzfn2o4qnlodg5yfuqtqzjxivo6vequroxwnpcmdokazmhs74/0x79bb3a0ee435f957ce4f54ee8c3cfadc7278da0c-rewards-369-331725.json",
		url)
}

func TestDownloadFile(t *testing.T) {
	const fileName = "rewards.json"
	fileContent := []byte(`{"Epoch":331725}`)

	t.Run("200 returns body", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(fileContent)
		}))
		defer srv.Close()

		client, err := pinata.NewClient(pinata.Config{Gateway: srv.URL + "/ipfs/%s/%s"})
		assert.Nil(t, err)

		content, err := client.DownloadFile("cid", fileName)
		assert.Nil(t, err)
		assert.Equal(t, fileContent, content)
	})

	t.Run("429 with Retry-After retries the same gateway once", func(t *testing.T) {
		var calls int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if atomic.AddInt32(&calls, 1) == 1 {
				w.Header().Set("Retry-After", "1")
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(fileContent)
		}))
		defer srv.Close()

		client, err := pinata.NewClient(pinata.Config{Gateway: srv.URL + "/ipfs/%s/%s"})
		assert.Nil(t, err)

		content, err := client.DownloadFile("cid", fileName)
		assert.Nil(t, err)
		assert.Equal(t, fileContent, content)
		assert.EqualValues(t, 2, atomic.LoadInt32(&calls))
	})

	t.Run("404 error contains rsp status err 404", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client, err := pinata.NewClient(pinata.Config{Gateway: srv.URL + "/ipfs/%s/%s"})
		assert.Nil(t, err)

		_, err = client.DownloadFile("cid", fileName)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "rsp status err 404")
	})

	t.Run("primary 500 then fallback 200 succeeds", func(t *testing.T) {
		primary := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer primary.Close()

		fallback := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(fileContent)
		}))
		defer fallback.Close()

		client, err := pinata.NewClient(pinata.Config{
			Gateway:          primary.URL + "/ipfs/%s/%s",
			FallbackGateways: []string{fallback.URL + "/ipfs/%s/%s"},
		})
		assert.Nil(t, err)

		content, err := client.DownloadFile("cid", fileName)
		assert.Nil(t, err)
		assert.Equal(t, fileContent, content)
	})

	t.Run("all gateways fail non-404 begins with primary status", func(t *testing.T) {
		primary := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer primary.Close()

		fallback := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
		}))
		defer fallback.Close()

		client, err := pinata.NewClient(pinata.Config{
			Gateway:          primary.URL + "/ipfs/%s/%s",
			FallbackGateways: []string{fallback.URL + "/ipfs/%s/%s"},
		})
		assert.Nil(t, err)

		_, err = client.DownloadFile("cid", fileName)
		assert.NotNil(t, err)
		assert.True(t, strings.HasPrefix(err.Error(), "rsp status err 500"), err.Error())
		assert.Contains(t, err.Error(), "rsp status err 502")
	})

	t.Run("gateway token appended as query parameter", func(t *testing.T) {
		var gotToken string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotToken = r.URL.Query().Get("pinataGatewayToken")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(fileContent)
		}))
		defer srv.Close()

		client, err := pinata.NewClient(pinata.Config{
			Gateway:      srv.URL + "/ipfs/%s/%s",
			GatewayToken: "secret-token",
		})
		assert.Nil(t, err)

		_, err = client.DownloadFile("cid", fileName)
		assert.Nil(t, err)
		assert.Equal(t, "secret-token", gotToken)
	})
}
