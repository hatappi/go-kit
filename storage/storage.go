package storage

import (
	"context"
	"fmt"

	"github.com/hatappi/go-kit/storage/provider"
)

type Storage interface {
	Save(ctx context.Context, filePath string, data []byte) (string, error)
	Get(ctx context.Context, filePath string) ([]byte, error)
	Ping(ctx context.Context) error
}

func NewStorage(serviceName string, conf *Config) (Storage, error) {
	switch conf.Type {
	case StorageTypeDisk:
		return provider.NewDisk(conf.Disk.RootDir), nil
	case StorageTypeS3:
		return provider.NewS3(conf.S3.BucketName, serviceName, conf.S3.Region)
	default:
		return nil, fmt.Errorf("invalid storage type: %s", conf.Type)
	}
}
