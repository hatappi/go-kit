package provider

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

type Disk struct {
	rootDir string
}

func NewDisk(root string) *Disk {
	return &Disk{
		rootDir: root,
	}
}

func (d *Disk) Save(filePath string, data []byte) (string, error) {
	savePath := path.Join(d.rootDir, filePath)

	if err := os.MkdirAll(filepath.Dir(savePath), 0755); err != nil {
		return "", err
	}

	file, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return "", err
	}

	return savePath, nil
}

func (d *Disk) Get(filePath string) ([]byte, error) {
	savePath := path.Join(d.rootDir, filePath)

	raw, err := ioutil.ReadFile(savePath)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

func (d *Disk) Ping() error {
	pingFilePath := path.Join(d.rootDir, "ping")

	fp, err := os.Create(pingFilePath)
	if err != nil {
		return err
	}

	if err := fp.Close(); err != nil {
		return err
	}

	if err := os.Remove(pingFilePath); err != nil {
		return err
	}

	return nil
}
