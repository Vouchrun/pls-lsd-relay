package nftstorage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"

	nftstorage "github.com/nftstorage/go-client"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/destorage"
)

var _ destorage.DeStorage = &NftStorage{}
var fileUrlFormatter string = "https://%s.ipfs.dweb.link/%s"

type NftStorage struct {
	apikey string
	log    *logrus.Entry
}

func (s *NftStorage) DownloadFile(cid string, fileName string) (content []byte, err error) {
	url := fmt.Sprintf(fileUrlFormatter, cid, fileName)
	s.log.WithFields(logrus.Fields{
		"url": url,
	}).Debug("DownloadFile")
	rsp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("rsp status err %d", rsp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	if len(bodyBytes) == 0 {
		return nil, fmt.Errorf("bodyBytes zero err")
	}
	return bodyBytes, nil
}

func (s *NftStorage) UploadFile(content []byte, path string) (cid string, err error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		return "", err
	}
	_, err = part.Write(content)
	if err != nil {
		return "", err
	}
	writer.Close()

	s.log.WithFields(logrus.Fields{
		"to":       "POST https://api.nft.storage/upload",
		"filename": path,
		"content":  string(content),
	}).Debug("UploadFile")

	r, _ := http.NewRequest("POST", "https://api.nft.storage/upload", body)
	r.Header.Add("Content-Type", writer.FormDataContentType())
	r.Header.Add("Authorization", "Bearer "+s.apikey)
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("upload error: status code[%s]", resp.Status)
	}

	respBody, _ := io.ReadAll(resp.Body)
	var data nftstorage.UploadResponse
	if err = json.Unmarshal(respBody, &data); err != nil {
		return "", err
	}

	if data.Value != nil && data.Value.Cid != nil {
		return *data.Value.Cid, nil
	}
	return "", fmt.Errorf("upload error: status code[%s] with body: %s", resp.Status, string(respBody))
}

func NewNftStorage(apikey string, logger *logrus.Entry) (*NftStorage, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger can not be nil")
	}
	return &NftStorage{
		apikey: apikey,
		log:    logger,
	}, nil
}
