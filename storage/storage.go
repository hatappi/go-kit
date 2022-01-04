package storage

import (
	"fmt"

	"github.com/hatappi/go-kit/storage/provider"
)

type Storage interface {
	Save(filePath string, data []byte) (string, error)
	Get(filePath string) ([]byte, error)
	Ping() error
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
