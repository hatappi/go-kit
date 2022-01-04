package storage

type StorageType string

const (
	StorageTypeDisk StorageType = "disk"
	StorageTypeS3   StorageType = "s3"
)
