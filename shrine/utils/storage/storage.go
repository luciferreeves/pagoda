package storage

import (
	"context"
	"errors"
	"io"
	"shrine/config"
	"shrine/messages"
	"shrine/utils/logger"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var Client *minio.Client

func init() {
	if config.Storage.AccessKey == "" || config.Storage.SecretKey == "" {
		logger.Infof("Storage", "MinIO credentials not configured, storage disabled")
		return
	}

	var err error
	Client, err = minio.New(config.Storage.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.Storage.AccessKey, config.Storage.SecretKey, ""),
		Secure: config.Storage.UseSSL,
	})
	if err != nil {
		logger.Fatalf("Storage", "Failed to initialize MinIO client: %v", err)
	}

	logger.Successf("Storage", "MinIO client initialized for %s", config.Storage.Endpoint)
}

func Upload(path string, reader io.Reader, size int64, contentType string) error {
	if Client == nil {
		return errors.New(messages.StorageNotConfigured)
	}
	_, err := Client.PutObject(context.Background(), config.Storage.Bucket, path, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func Delete(path string) error {
	if Client == nil {
		return errors.New(messages.StorageNotConfigured)
	}
	return Client.RemoveObject(context.Background(), config.Storage.Bucket, path, minio.RemoveObjectOptions{})
}

func ResolveCDN(path string) string {
	if path == "" {
		return ""
	}
	return strings.TrimRight(config.Storage.CDN, "/") + "/" + config.Storage.Bucket + "/" + path
}

func PathFromCDN(cdnURL string) string {
	prefix := strings.TrimRight(config.Storage.CDN, "/") + "/" + config.Storage.Bucket + "/"
	if !strings.HasPrefix(cdnURL, prefix) {
		return ""
	}
	return strings.TrimPrefix(cdnURL, prefix)
}