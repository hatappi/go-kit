package storage

import "testing"

func TestNewStorage(t *testing.T) {
	testCases := []struct {
		name   string
		config *Config

		wantErr bool
	}{
		{
			name: "storage type is disk",
			config: &Config{
				Type: "disk",
				Disk: struct {
					RootDir string `envconfig:"ROOT_DIR"`
				}{
					RootDir: "/",
				},
			},
			wantErr: false,
		},
		{
			name: "storage type is s3",
			config: &Config{
				Type: "s3",
				S3: struct {
					BucketName string `envconfig:"BUCKET_NAME"`
					Region     string `default:"ap-northeast-1" envconfig:"REGION"`
				}{
					BucketName: "test",
				},
			},
			wantErr: false,
		},
		{
			name: "storage type is invalid",
			config: &Config{
				Type: "test",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			_, err := NewStorage("test", tc.config)
			if (err != nil) != tc.wantErr {
				t.Fatalf("err: %v", err)
			}
		})
	}
}
