package storage

import (
	"context"
	"fmt"
	"io"
	"shrine/config"
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
		return fmt.Errorf("storage not configured")
	}
	_, err := Client.PutObject(context.Background(), config.Storage.Bucket, path, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func Delete(path string) error {
	if Client == nil {
		return fmt.Errorf("storage not configured")
	}
	return Client.RemoveObject(context.Background(), config.Storage.Bucket, path, minio.RemoveObjectOptions{})
}

func ResolveCDN(path string) string {
	if path == "" {
		return ""
	}
	return strings.TrimRight(config.Storage.CDN, "/") + "/" + config.Storage.Bucket + "/" + path
}