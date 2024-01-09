package destorage

type DeStorage interface {
	DownloadFile(cid, fileName string) (content []byte, err error)
	// Upload a file to Decentralized Storage and get the CID for it
	UploadFile(content []byte, path string) (cid string, err error)
}
