package uproot

import (
	"embed"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
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

func New(eFS *embed.FS) *uproot {
	uproot := &uproot{
		eFS: eFS,
	}

	err := uproot.scanFS()
	if err != nil {
		log.Fatalln(err)
	}
	uproot.setTmpDir()

	return uproot
}

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

func (u *uproot) CopyToTmpDir() bool {

	for _, file := range u.files {
		content, err := u.eFS.ReadFile(file)
		if err != nil {
			log.Fatalln(err)
		}

		filename := filepath.Join(u.tmpDir, file)
		os.MkdirAll(filepath.Dir(filename), 0666)
		if err := os.WriteFile(filename, content, 0666); err != nil {
			log.Fatalln(err)
		}
	}

	return true
}

func (u *uproot) setTmpDir() {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		log.Printf("Failed to read build info")

	}
	u.tmpDir = filepath.Join(os.TempDir(), bi.Main.Path)
}

func (u *uproot) GetTmpDir() string {
	return u.tmpDir
}
func (u *uproot) GetFS() *embed.FS {
	return u.eFS
}
func (u *uproot) GetFiles() []string {
	return u.files
}
