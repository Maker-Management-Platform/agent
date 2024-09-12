package libfs

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"path/filepath"

	"github.com/eduardooliveira/stLib/v2/config"
	"github.com/mholt/archiver/v4"
)

type bundleFS struct {
	FS
	bundleLocation string
}

func newBundleFS(parentFS LibFS, path string) (LibFS, error) {

	base := filepath.Base(path)
	rtn := FS{
		FileSystem: config.FileSystem{
			Name: parentFS.GetName(),
			Path: path,
			Kind: "bundle",
		},
		discovarable: true,
	}

	var err error
	switch filepath.Ext(base) {
	case ".zip", ".rar", ".7z", ".tar", ".3mf":
		rtn.FS, err = archiver.FileSystem(context.Background(), filepath.Join(parentFS.GetLocation(), path))
		if err != nil {
			return nil, err
		}
	}
	return bundleFS{
		FS:             rtn,
		bundleLocation: filepath.Join(parentFS.GetLocation(), path),
	}, nil
}

func (fs bundleFS) GetFS() fs.FS {
	return fs.FS
}
func (fs bundleFS) GetName() string {
	return fs.Name
}
func (fs bundleFS) GetLocation() string {
	return fs.Path
}
func (fs bundleFS) Kind() string {
	return "bundle"
}

func (fs bundleFS) Writable() bool {
	return false
}

func (fs bundleFS) Create(name string) (file io.WriteCloser, err error) {
	return nil, errors.New("write not supported")
}

func (fs bundleFS) Mkdir(name string) error {
	return errors.New("mkdir not supported")
}

func (fs bundleFS) setDiscovarable(d bool) {
	fs.discovarable = d
}
