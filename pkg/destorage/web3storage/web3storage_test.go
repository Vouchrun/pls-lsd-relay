package web3storage_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stafiprotocol/eth-lsd-relay/pkg/destorage/web3storage"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	priv := os.Getenv("WEB3STORAGE_PRIVATE_KEY")                // MgCYR...
	proofFile := os.Getenv("WEB3STORAGE_DELEGATION_PROOF_FILE") // /somewhere/proof.ucan
	spaceDid := os.Getenv("WEB3STORAGE_SPACE_DID")              // did:key:abcDEF

	_, err := web3storage.NewStorage(
		proofFile,
		spaceDid,
		priv,
	)
	assert.Nil(t, err)
}

func TestUploadAndDownload(t *testing.T) {
	priv := os.Getenv("WEB3STORAGE_PRIVATE_KEY")                // MgCYR...
	proofFile := os.Getenv("WEB3STORAGE_DELEGATION_PROOF_FILE") // /somewhere/proof.ucan
	spaceDid := os.Getenv("WEB3STORAGE_SPACE_DID")              // did:key:abcDEF

	s, err := web3storage.NewStorage(
		proofFile,
		spaceDid,
		priv,
	)
	assert.Nil(t, err)

	fileName := "got-20240625-3.txt"
	fileContent := []byte("Valar Morghulis")
	cid, err := s.UploadFile(fileContent, fileName)
	assert.Nil(t, err)
	assert.NotEmpty(t, cid)
	fmt.Println("file name:", fileName)
	fmt.Println("cid:", cid)

	downloadContent, err := s.DownloadFile(cid, fileName)
	assert.Nil(t, err)
	assert.Equal(t, string(fileContent), string(downloadContent))
	fmt.Println("content:", string(downloadContent))
}
