package provider

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDiskSave(t *testing.T) {
	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	diskProvider := &Disk{
		rootDir: dir,
	}

	ctx := context.Background()
	savedPath, err := diskProvider.Save(ctx, "test.txt", []byte("test"))
	if err != nil {
		t.Fatal(err)
	}

	actual, err := ioutil.ReadFile(savedPath)
	if err != nil {
		t.Fatal(err)
	}

	expected := []byte("test")
	if d := cmp.Diff(expected, actual); d != "" {
		t.Fatalf("unexpected contents. %s", d)
	}
}

func TestDiskGet(t *testing.T) {
	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	file, err := os.Create(path.Join(dir, "test.txt"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write([]byte("test"))
	if err != nil {
		t.Fatal(err)
	}

	diskProvider := &Disk{
		rootDir: dir,
	}

	ctx := context.Background()
	actual, err := diskProvider.Get(ctx, "test.txt")
	if err != nil {
		t.Fatal(err)
	}

	expected := []byte("test")
	if d := cmp.Diff(expected, actual); d != "" {
		t.Fatalf("unexpected contents. %s", d)
	}
}

func TestDiskDelete(t *testing.T) {
	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	file, err := os.Create(path.Join(dir, "test.txt"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	diskProvider := &Disk{
		rootDir: dir,
	}

	ctx := context.Background()
	err = diskProvider.Delete(ctx, "test.txt")
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(path.Join(dir, "test.txt"))
	if os.IsExist(err) {
		t.Fatal("file was not deleted")
	}
}

func TestDiskPing(t *testing.T) {
	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	diskProvider := &Disk{
		rootDir: dir,
	}

	ctx := context.Background()
	err = diskProvider.Ping(ctx)
	if err != nil {
		t.Fatal(err)
	}
}
