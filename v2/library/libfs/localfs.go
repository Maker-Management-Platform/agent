package libfs

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/eduardooliveira/stLib/v2/config"
	"github.com/eduardooliveira/stLib/v2/utils"
)

type localFS struct {
	lFS
}

func newLocalFS(cfgFS config.FileSystem) LibFS {
	return &localFS{
		lFS: lFS{
			FileSystem:   cfgFS,
			FS:           os.DirFS(cfgFS.Path),
			discovarable: true,
		},
	}
}

func (fs localFS) GetFS() fs.FS {
	return fs.lFS
}
func (fs localFS) GetName() string {
	return fs.lFS.Name
}
func (fs localFS) GetRoot() string {
	return fs.Path
}
func (fs localFS) GetLocation() string {
	return fs.Path
}
func (fs localFS) Kind() string {
	return "local"
}

func (fs localFS) Writable() bool {
	return true
}

func (fs localFS) Create(name string) (file io.WriteCloser, err error) {
	if filepath.Base(name) != "." {
		err := utils.CreateFolder(filepath.Join(fs.Path, filepath.Dir(name)))
		if err != nil {
			return nil, err
		}
	}
	return os.Create(filepath.Join(fs.Path, name))
}

func (fs localFS) Mkdir(name string) error {
	return os.MkdirAll(filepath.Join(fs.Path, name), 0755)
}

func (fs *localFS) setDiscovarable(d bool) {
	fs.discovarable = d
}

func (fs localFS) Remove(name string) error {
	return os.RemoveAll(filepath.Join(fs.Path, name))
}
