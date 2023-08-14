// Package uproot provides simplified access to files that are embedded in an embed.FS.
//
// No support for files which are ambedded as []byte or string.
package uproot

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
)

type Uproot interface {
	FS() *embed.FS
	TmpDir() string
	Files() []string
}

type uproot struct {
	eFS    *embed.FS
	files  []string
	tmpDir string
}

const tmpDirName string = "uproot-fs"

/*Exported functions*/

// Creates a new instance of uproot
func New(eFS *embed.FS) (*uproot, error) {
	u := &uproot{
		eFS:    eFS,
		tmpDir: filepath.Join(os.TempDir(), tmpDirName),
	}
	err := u.scanFS()
	if err != nil {
		return nil, err
	}
	return u, nil
}

// Copies all files to the given directory.
//
// If the directory does not exist it will be created
func (u *uproot) CopyTo(dir string) error {
	return u.copy(dir)
}

// Copies all files of the file system to the temp directory
func (u *uproot) CopyToTmp() error {
	return u.copy(u.tmpDir)
}

// Deletes the temp directory including the contents
func (u *uproot) RemoveTmp() error {
	err := os.RemoveAll(u.tmpDir)
	if err != nil {
		return err
	}
	return nil
}

// Returns the path to the temp directory
func (u *uproot) TmpDir() string {
	return u.tmpDir
}

// Returns the pointer to embedded FS
func (u *uproot) FS() *embed.FS {
	return u.eFS
}

// Returns the scanned files
func (u *uproot) Files() []string {
	return u.files
}

/*Unexported functions*/

// Scans the directory and adds all files in uproot.files
func (u *uproot) scanFS() error {
	err := fs.WalkDir(u.eFS, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		u.files = append(u.files, path)

		return nil
	})
	return err
}

func (u *uproot) copy(dir string) error {
	for _, file := range u.files {
		content, err := u.eFS.ReadFile(file)
		if err != nil {
			return err
		}
		path := filepath.Join(dir, filepath.Dir(file))
		filepath := filepath.Join(path, filepath.Base(file))
		err = os.MkdirAll(path, 0777)
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath, content, 0777); err != nil {
			return err
		}
	}
	return nil
}
