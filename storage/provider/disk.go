package provider

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/hatappi/go-kit/storage/option"
)

type Disk struct {
	rootDir string
}

func NewDisk(root string) *Disk {
	return &Disk{
		rootDir: root,
	}
}

func (d *Disk) Save(ctx context.Context, filePath string, data []byte, opts ...option.SaveOptionFunc) (string, error) {
	var saveOpt option.SaveOption
	for _, opt := range opts {
		opt(&saveOpt)
	}

	fullPath := d.fileFullPath(filePath)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return "", err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return "", err
	}

	return fullPath, nil
}

func (d *Disk) Get(ctx context.Context, filePath string) ([]byte, error) {
	fullPath := d.fileFullPath(filePath)

	_, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return nil, nil
	}

	raw, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

func (d *Disk) Delete(ctx context.Context, filePath string) error {
	if err := os.Remove(d.fileFullPath(filePath)); err != nil {
		return err
	}

	return nil
}

func (d *Disk) Ping(ctx context.Context) error {
	filePath := "ping"

	_, err := d.Save(ctx, filePath, []byte("test"))
	if err != nil {
		return err
	}

	b, err := d.Get(ctx, filePath)
	if err != nil {
		return err
	}

	if b == nil {
		return errors.New("file does not exist")
	}

	if err := d.Delete(ctx, filePath); err != nil {
		return err
	}

	return nil
}

func (d *Disk) fileFullPath(filePath string) string {
	return path.Join(d.rootDir, filePath)
}
