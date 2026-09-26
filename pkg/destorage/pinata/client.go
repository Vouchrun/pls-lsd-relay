package pinata

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/stafiprotocol/eth-lsd-relay/pkg/destorage"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/utils"
)

var _ destorage.DeStorage = &Client{}

type Client struct {
	endpoint     string
	apikey       string
	gateways     []string
	gatewayToken string
	httpClient   *http.Client
}

type Config struct {
	Endpoint         string
	Apikey           string
	Gateway          string   // download URL template, two %s (cid, filename)
	FallbackGateways []string // ordered fallback templates, same format
	GatewayToken     string   // optional token for a dedicated/private gateway
}

const (
	defaultEndpoint = "https://api.pinata.cloud"
	DefaultGateway  = "https://gateway.pinata.cloud/ipfs/%s/%s"
)

func NewClient(cfg Config) (*Client, error) {
	endpoint := cfg.Endpoint
	if endpoint == "" {
		endpoint = defaultEndpoint
	}

	gateways := make([]string, 0, 1+len(cfg.FallbackGateways))
	if cfg.Gateway != "" {
		gateways = append(gateways, cfg.Gateway)
	} else {
		gateways = append(gateways, DefaultGateway)
	}
	for _, gateway := range cfg.FallbackGateways {
		if gateway != "" {
			gateways = append(gateways, gateway)
		}
	}

	c := &Client{
		endpoint:     endpoint,
		apikey:       cfg.Apikey,
		gateways:     gateways,
		gatewayToken: cfg.GatewayToken,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
	}

	return c, nil
}

func (c *Client) StartUnpinFiles(pinDur time.Duration) {
	if pinDur <= 0 {
		return
	}

	utils.SafeGoWithRestart(func() {
		for {
			count, err := c.UnpinFilesCreatedBefore(time.Now().Add(-pinDur))
			if err != nil {
				slog.Error("[pinata]: fail to unpin outdated files", "err", err)
			}
			if count > 0 {
				slog.Info("[pinata]: successfully unpinned outdated files", "count", count)
			}
			time.Sleep(time.Hour * 24)
		}
	})
}

func (c *Client) DownloadFile(cid, fileName string) (content []byte, err error) {
	outcomes := make([]string, 0, len(c.gateways))
	var primaryStatus int
	var primaryErr error
	allNotFound := true

	for i, tmpl := range c.gateways {
		bodyBytes, statusCode, gatewayErr := c.downloadFromGateway(tmpl, cid, fileName)
		if gatewayErr == nil {
			return bodyBytes, nil
		}

		outcomes = append(outcomes, fmt.Sprintf("%s: %v", tmpl, gatewayErr))

		if i == 0 {
			if statusCode == 0 {
				primaryErr = gatewayErr
			} else {
				primaryStatus = statusCode
			}
		}
		if statusCode != http.StatusNotFound {
			allNotFound = false
		}
	}

	// primary gateway failed at transport level: preserve the verbatim error as before
	if primaryErr != nil {
		return nil, primaryErr
	}

	statusCode := primaryStatus
	if allNotFound {
		statusCode = http.StatusNotFound
	}

	return nil, fmt.Errorf("rsp status err %d (%s)", statusCode, strings.Join(outcomes, "; "))
}

func (c *Client) downloadFromGateway(tmpl, cid, fileName string) (content []byte, statusCode int, err error) {
	downloadURL := fmt.Sprintf(tmpl, cid, fileName)
	if c.gatewayToken != "" {
		parsed, err := url.Parse(downloadURL)
		if err != nil {
			return nil, 0, err
		}
		query := parsed.Query()
		query.Set("pinataGatewayToken", c.gatewayToken)
		parsed.RawQuery = query.Encode()
		downloadURL = parsed.String()
	}

	for attempt := 0; ; attempt++ {
		rsp, err := c.httpClient.Get(downloadURL)
		if err != nil {
			return nil, 0, err
		}

		if rsp.StatusCode == http.StatusOK {
			bodyBytes, err := io.ReadAll(rsp.Body)
			rsp.Body.Close()
			if err != nil {
				return nil, rsp.StatusCode, err
			}
			if len(bodyBytes) == 0 {
				return nil, rsp.StatusCode, fmt.Errorf("bodyBytes zero err")
			}
			return bodyBytes, rsp.StatusCode, nil
		}

		retryable := rsp.StatusCode == http.StatusTooManyRequests && attempt == 0
		retryAfter := 0
		if retryable {
			retryAfter, _ = strconv.Atoi(strings.TrimSpace(rsp.Header.Get("Retry-After")))
		}
		rsp.Body.Close()

		if retryable && retryAfter > 0 && retryAfter <= 5 {
			time.Sleep(time.Duration(retryAfter) * time.Second)
			continue
		}

		return nil, rsp.StatusCode, fmt.Errorf("rsp status err %d", rsp.StatusCode)
	}
}

func (c *Client) UploadFile(content []byte, path string) (cid string, err error) {
	tempDir, err := os.MkdirTemp("", "ethlsd")
	if err != nil {
		return "", fmt.Errorf("create temp dir error: %w", err)
	}
	filepath := joinPath(tempDir, filepath.Base(path))
	if err = os.WriteFile(filepath, content, 0600); err != nil {
		return "", fmt.Errorf("write to file error: %w", err)
	}
	defer os.RemoveAll(tempDir)

	return c.uploadFile(tempDir)
}

type UploadResponse struct {
	IpfsHash    string `json:"IpfsHash"`
	PinSize     int    `json:"PinSize"`
	Timestamp   string `json:"Timestamp"`
	IsDuplicate bool   `json:"isDuplicate"`
}

type Options struct {
	CidVersion int `json:"cidVersion"`
}

type Metadata struct {
	Name      string                 `json:"name"`
	KeyValues map[string]interface{} `json:"keyvalues"`
}

func (c *Client) uploadFile(filePath string) (string, error) {
	stats, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		fmt.Println("File or folder does not exist")
		return "", errors.Join(err, errors.New("file or folder does not exist"))
	}

	files, err := pathsFinder(filePath, stats)
	if err != nil {
		return "", err
	}

	body := &bytes.Buffer{}
	contentType, err := createMultipartRequest(filePath, files, body, stats, 1)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/pinning/pinFileToIPFS", c.endpoint)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", errors.Join(err, errors.New("failed to create the request"))
	}
	req.Header.Set("content-type", contentType)

	var response UploadResponse
	if err = c.do(req, &response); err != nil {
		return "", err
	}

	return response.IpfsHash, nil
}

func createMultipartRequest(filePath string, files []string, body io.Writer, stats os.FileInfo, version int) (string, error) {
	contentType := ""
	writer := multipart.NewWriter(body)

	fileIsASingleFile := !stats.IsDir()
	for _, f := range files {
		file, err := os.Open(f)
		if err != nil {
			return contentType, err
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Fatal("could not close file")
			}
		}(file)

		var part io.Writer
		if fileIsASingleFile {
			part, err = writer.CreateFormFile("file", filepath.Base(f))
		} else {
			relPath, _ := filepath.Rel(filePath, f)
			part, err = writer.CreateFormFile("file", filepath.Join(stats.Name(), relPath))
		}
		if err != nil {
			return contentType, err
		}
		_, err = io.Copy(part, file)
		if err != nil {
			return contentType, err
		}
	}

	pinataOptions := Options{
		CidVersion: version,
	}

	optionsBytes, err := json.Marshal(pinataOptions)
	if err != nil {
		return contentType, err
	}
	err = writer.WriteField("pinataOptions", string(optionsBytes))

	if err != nil {
		return contentType, err
	}

	pinataMetadata := Metadata{
		Name: stats.Name(),
		KeyValues: map[string]any{
			"created_at": time.Now().Unix(),
		},
	}
	metadataBytes, err := json.Marshal(pinataMetadata)
	if err != nil {
		return contentType, err
	}
	_ = writer.WriteField("pinataMetadata", string(metadataBytes))
	err = writer.Close()
	if err != nil {
		return contentType, err
	}

	contentType = writer.FormDataContentType()

	return contentType, nil
}

type Pin struct {
	Id            string   `json:"id"`
	IPFSPinHash   string   `json:"ipfs_pin_hash"`
	Size          int      `json:"size"`
	UserId        string   `json:"user_id"`
	DatePinned    string   `json:"date_pinned"`
	DateUnpinned  *string  `json:"date_unpinned"`
	Metadata      Metadata `json:"metadata"`
	MimeType      string   `json:"mime_type"`
	NumberOfFiles int      `json:"number_of_files"`
}

type ListResponse struct {
	Rows []Pin `json:"rows"`
}

func (c *Client) UnpinFilesCreatedBefore(before time.Time) (count int, err error) {
	query := fmt.Sprintf(`status=pinned&metadata[keyvalues][created_at]={"value":"%d","op":"lt"}`, before.Unix())
	resp, err := c.listFiles(query)
	if err != nil {
		return 0, err
	}
	for _, row := range resp.Rows {
		if err = c.Delete(row.IPFSPinHash); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

func (c *Client) listFiles(query string) (ListResponse, error) {
	url := fmt.Sprintf("%s/data/pinList?%s", c.endpoint, query)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ListResponse{}, errors.Join(err, errors.New("failed to create the request"))
	}
	req.Header.Set("content-type", "application/json")
	var response ListResponse
	if err = c.do(req, &response); err != nil {
		return ListResponse{}, err
	}

	return response, nil
}

func (c *Client) Delete(cid string) error {
	url := fmt.Sprintf("%s/pinning/unpin/%s", c.endpoint, cid)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("fail to create request: %w", err)
	}

	return c.do(req, nil)
}

func (c *Client) do(req *http.Request, res any) error {
	req.Header.Set("Authorization", "Bearer "+c.apikey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("fail to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("server returned an error: %d", resp.StatusCode)
	}

	if res != nil {
		if err = json.NewDecoder(resp.Body).Decode(res); err != nil {
			return fmt.Errorf("fail to decode response: %w", err)
		}
	}
	return nil
}

func pathsFinder(filePath string, stats os.FileInfo) ([]string, error) {
	var err error
	files := make([]string, 0)
	fileIsASingleFile := !stats.IsDir()
	if fileIsASingleFile {
		files = append(files, filePath)
		return files, err
	}
	err = filepath.Walk(filePath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		})

	if err != nil {
		return nil, err
	}

	return files, err
}

func joinPath(dir, name string) string {
	if len(dir) > 0 && os.IsPathSeparator(dir[len(dir)-1]) {
		return dir + name
	}
	return dir + string(os.PathSeparator) + name
}
