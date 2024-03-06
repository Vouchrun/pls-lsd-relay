package local_store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Info struct {
	SyncedHeight uint64
	Address      string `json:"-"`
}

type LocalStore struct {
	mu   sync.Mutex
	path string
}

func NewLocalStore(path string) (*LocalStore, error) {
	s := LocalStore{
		path: path,
	}

	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open or create local store file err: %w", err)
	}
	defer f.Close()

	return &s, nil
}

func (s *LocalStore) Read(addr string) (*Info, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	content, err := s.readContent()
	if err != nil {
		return nil, err
	}
	info, ok := content[addr]
	if !ok {
		return nil, nil // address info does not exist
	}
	info.Address = addr
	return &info, nil
}

func (s *LocalStore) Update(update Info) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	content, err := s.readContent()
	if err != nil {
		return err
	}
	content[update.Address] = update

	bytes, err := json.Marshal(content)
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, bytes, 0644)
}

func (s *LocalStore) readContent() (map[string]Info, error) {
	content, err := os.ReadFile(s.path)
	if err != nil {
		return nil, err
	}
	content = bytes.TrimSpace(content)
	info := map[string]Info{}
	if len(content) == 0 {
		return info, nil
	}
	if err = json.Unmarshal(content, &info); err != nil {
		return nil, err
	}
	return info, nil
}
