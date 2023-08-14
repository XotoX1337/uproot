package uproot

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed test-assets/*
var assets embed.FS

func TestCopyTo(t *testing.T) {
	assert := assert.New(t)
	u, err := New(&assets)
	assert.IsType(nil, err, "error should be nil")
	path := "test-copy"
	u.CopyTo(path)
	assert.DirExistsf(path, "dir %s should exist", path)

	for _, file := range u.Files() {
		assert.FileExistsf(filepath.Join(path, file), "file %s should exist", file)
	}
	os.RemoveAll(path)
	assert.NoDirExistsf(path, "dir %s should not exist", path)
}

func TestFS(t *testing.T) {
	assert := assert.New(t)
	u, err := New(&assets)
	assert.IsType(nil, err, "error should be nil")
	assert.IsType(&embed.FS{}, u.FS(), "should be type embed.FS")
}

func TestFiles(t *testing.T) {
	assert := assert.New(t)
	u, err := New(&assets)
	assert.IsType(nil, err, "error should be nil")
	assert.IsType([]string{}, u.Files(), "should be type []string")

	fs.WalkDir(u.eFS, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			assert.DirExists(path, "dir should exist")
		} else {
			assert.FileExists(path, "file should exist")
		}
		return nil
	})
}

func TestTmpDir(t *testing.T) {
	assert := assert.New(t)
	u, err := New(&assets)
	assert.IsType(nil, err, "error should be nil")
	assert.Equal(filepath.Join(os.TempDir(), "uproot-fs"), u.TmpDir())
}

func TestRemoveTmp(t *testing.T) {
	assert := assert.New(t)
	u, err := New(&assets)
	assert.IsType(nil, err, "error should be nil")
	err = u.CopyToTmp()
	assert.IsType(nil, err, "error while copying to temp dir")
	assert.DirExists(u.TmpDir(), "temp dir should exist")
	err = u.RemoveTmp()
	assert.IsType(nil, err, "error while deleting temp dir")
	assert.NoDirExists(u.TmpDir(), "temp dir should not exist")

}
