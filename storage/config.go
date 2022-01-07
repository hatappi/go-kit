package storage

type Config struct {
	Type StorageType `envconfig:"TYPE" validate:"oneof=disk s3"`

	Disk struct {
		RootDir string `envconfig:"ROOT_DIR"`
	} `envconfig:"DISK"`

	S3 struct {
		BucketName string `envconfig:"BUCKET_NAME"`
		Region     string `default:"ap-northeast-1" envconfig:"REGION"`
	} `envconfig:"S3"`
}
