package storage_service

import (
	"log"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	BasePath string
}

func (s *LocalStorage) Save(filename string, data []byte) (string, error) {
	if _, err := os.Stat(s.BasePath); os.IsNotExist(err) {
		_ = os.Mkdir(s.BasePath, 0755)
	}

	path := filepath.Join(s.BasePath, filename)
	err := os.WriteFile(path, data, 0755)

	if err != nil {
		log.Printf("error in save file: %s", err)
		return "", err
	}

	return path, nil
}

func (s *LocalStorage) Delete(filename string) error {
	path := filepath.Join(s.BasePath, filename)
	err := os.Remove(path)

	if err != nil {
		log.Printf("error in delete file: %s", err)
	}

	return err
}
