package nftstorage_test

import (
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/destorage/nftstorage"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	return
	_, err := nftstorage.NewNftStorage("", logrus.NewEntry(logrus.New()))
	assert.Equal(t, fmt.Errorf("apikey must be set"), err)
	storage, err := nftstorage.NewNftStorage("abc", logrus.NewEntry(logrus.New()))
	assert.Nil(t, err)
	assert.NotNil(t, storage)
}

func TestUpload(t *testing.T) {
	return
	apikey := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJkaWQ6ZXRocjoweDY4MzIyNTIzMTUxODcxQkYzZDZmMGEwM2YwMDNGYkYxNGQ1MjA4N2MiLCJpc3MiOiJuZnQtc3RvcmFnZSIsImlhdCI6MTcwNDc4MTg3NjE1OSwibmFtZSI6InRlc3QgcHJvZ3JhbSJ9.4WqIoZDH8nof_ypF6klfA0LW6OT4baMkwMv6tqL3bbA"
	s, err := nftstorage.NewNftStorage(apikey, logrus.NewEntry(logrus.New()))
	assert.Nil(t, err)
	filePath := "hello-nft.storage.txt"
	content := []byte("hello nft.storage")
	cid, err := s.UploadFile(content, filePath)
	assert.Nil(t, err)
	assert.NotEqual(t, "", cid)
	fmt.Println("cid:", cid)
}

func TestDownload(t *testing.T) {
	return
	apikey := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJkaWQ6ZXRocjoweDY4MzIyNTIzMTUxODcxQkYzZDZmMGEwM2YwMDNGYkYxNGQ1MjA4N2MiLCJpc3MiOiJuZnQtc3RvcmFnZSIsImlhdCI6MTcwNDc4MTg3NjE1OSwibmFtZSI6InRlc3QgcHJvZ3JhbSJ9.4WqIoZDH8nof_ypF6klfA0LW6OT4baMkwMv6tqL3bbA"
	s, err := nftstorage.NewNftStorage(apikey, logrus.NewEntry(logrus.New()))
	assert.Nil(t, err)
	filePath := "hello-nft.storage.txt"
	cid := "bafybeibc5bswvtk4geg746wzqtticmdalcegwuhrhm2k6hg72dwgeczjj4"
	content, err := s.DownloadFile(cid, filePath)
	assert.Nil(t, err)
	expectContent := []byte("hello nft.storage")
	assert.Equal(t, expectContent, content)
}
