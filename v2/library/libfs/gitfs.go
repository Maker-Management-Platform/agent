package libfs

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/url"

	"github.com/eduardooliveira/stLib/v2/config"
	gfs "github.com/hairyhenderson/go-fsimpl/gitfs"
)

type gitfs struct {
	FS
	config gitFSConfig
}
type gitFSConfig struct {
	url string
}

func newGitFS(cfgFS config.FileSystem) (LibFS, error) {
	rtn := gitfs{
		FS: FS{
			FileSystem:   cfgFS,
			discovarable: true,
		},
	}

	var base *url.URL
	var err error
	if u, ok := cfgFS.Config["url"].(string); ok {
		rtn.config.url = u
		base, err = url.Parse(rtn.config.url)
		if err != nil {
			return nil, fmt.Errorf("invalid gitfs url: %s", err.Error())
		}
	} else {
		return nil, errors.New("gitfs requires a url")
	}

	fsys, err := gfs.New(base)
	if err != nil {
		slog.Error(err.Error())
	}
	rtn.FS.FS = fsys

	return rtn, nil
}

func (fs gitfs) GetFS() fs.FS {
	return fs.FS
}
func (fs gitfs) GetName() string {
	return fs.Name
}
func (fs gitfs) GetLocation() string {
	return fs.Path
}
func (fs gitfs) Kind() string {
	return "gitfs"
}

func (fs gitfs) Writable() bool {
	return false
}

func (fs gitfs) Create(name string) (file io.WriteCloser, err error) {
	return nil, errors.New("gitfs is read-only")
}

func (fs gitfs) Mkdir(name string) error {
	return errors.New("gitfs is read-only")
}

func (fs gitfs) setDiscovarable(d bool) {
	fs.discovarable = d
}
