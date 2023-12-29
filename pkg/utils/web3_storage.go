package utils

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/web3-storage/go-w3s-client"
)

var nodeRewardsFileUrlRaw string = "https://%s.ipfs.dweb.link/%s"
var nodeRewardsFileNameRaw string = "%s-nodeRewards-%d.json"

func NodeRewardsFileNameAtEpoch(lsdToken string, epoch uint64) string {
	return fmt.Sprintf(nodeRewardsFileNameRaw, strings.ToLower(lsdToken), epoch)
}

func DownloadWeb3File(cid, fileName string) ([]byte, error) {
	url := fmt.Sprintf(nodeRewardsFileUrlRaw, cid, fileName)
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

// Compress and upload a file to Web3.Storage and get the CID for it
func UploadFileToWeb3Storage(client w3s.Client, wrapperBytes []byte, compressedPath string) (string, error) {

	// todo compress file
	compressedBytes := wrapperBytes

	// Create the compressed tree file
	compressedFile, err := os.Create(compressedPath)
	if err != nil {
		return "", fmt.Errorf("error creating file [%s]: %w", compressedPath, err)
	}
	defer compressedFile.Close()

	// Write the compressed data to the file
	_, err = compressedFile.Write(compressedBytes)
	if err != nil {
		return "", fmt.Errorf("error writing to %s: %w", compressedPath, err)
	}

	// Rewind it to the start
	_, err = compressedFile.Seek(0, 0)
	if err != nil {
		return "", fmt.Errorf("seek failed: %w", err)
	}

	// Upload it
	cid, err := client.Put(context.Background(), compressedFile)
	if err != nil {
		return "", fmt.Errorf("error uploading: %w", err)
	}

	return cid.String(), nil

}
