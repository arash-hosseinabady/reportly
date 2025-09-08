package storage_service

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type StorageService interface {
	Save(fileName string, data []byte) (string, error)
	Delete(fileName string) error
}

const (
	localStorage = "local"
	minioStorage = "minio"
)

var Driver string

func NewStorageService() (StorageService, error) {
	Driver = os.Getenv("STORAGE_DRIVER")

	switch Driver {
	case localStorage:
		path := os.Getenv("LOCAL_STORAGE_PATH")
		return &LocalStorage{BasePath: path}, nil
	case minioStorage:
		endpoint := os.Getenv("MINIO_ENDPOINT")
		accessKey := os.Getenv("MINIO_ACCESS_KEY")
		secretKey := os.Getenv("MINIO_SECRET_KEY")
		bucket := os.Getenv("MINIO_BUCKET")
		useSSL, _ := strconv.ParseBool(os.Getenv("MINIO_USE_SSL"))

		return NewMinioStorage(endpoint, accessKey, secretKey, bucket, useSSL)
	default:
		log.Printf("unsupported storage driver: %s", Driver)
		return nil, fmt.Errorf("unsupported storage driver: %s", Driver)
	}
}
