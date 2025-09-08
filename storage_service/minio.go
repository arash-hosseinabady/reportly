package storage_service

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

type MinioStorage struct {
	Client     *minio.Client
	BucketName string
}

func NewMinioStorage(endpoint, accessKey, secretKey, bucketName string, useSSL bool) (*MinioStorage, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		log.Printf("error creating minio client: %s", err)
		return nil, err
	}

	return &MinioStorage{Client: minioClient, BucketName: bucketName}, nil
}

func (s *MinioStorage) Save(fileName string, data []byte) (string, error) {
	_, err := s.Client.PutObject(
		context.Background(),
		s.BucketName,
		fileName,
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{ContentType: "application/octet-stream"},
	)

	if err != nil {
		log.Printf("error uploading file: %s", err)
		return "", err
	}

	return fileName, nil
}

func (s *MinioStorage) Delete(fileName string) error {
	err := s.Client.RemoveObject(context.Background(), s.BucketName, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		log.Printf("error deleting file: %s", err)
	}
	return err
}
