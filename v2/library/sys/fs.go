package sys

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/mholt/archiver/v4"
)

type FS struct {
	fs.FS
	Name string `json:"name"`
	Kind string `json:"kind"`
	Path string `json:"path"`
}

func GetBundleFS(parentFS FS, path string) (FS, error) {
	if parentFS.Kind != "local" {
		return FS{}, errors.New("only local FS is supported")
	}
	base := filepath.Base(path)
	rtn := FS{
		FS:   parentFS,
		Name: base,
		Path: filepath.Join(parentFS.Path, path),
		Kind: "bundle",
	}

	var err error
	switch filepath.Ext(base) {
	case ".zip", ".rar", ".7z", ".tar", ".3mf":
		rtn.FS, err = archiver.FileSystem(context.Background(), filepath.Join(parentFS.Path, path))
		if err != nil {
			return FS{}, err
		}
	}

	return rtn, nil
}

func IsBundle(path string) bool {
	switch filepath.Ext(path) {
	case ".zip", ".rar", ".7z", ".tar", ".3mf":
		return true
	}
	return false
}

func GetFS(kind, name, path string) (FS, error) {
	var err error
	f := FS{
		Name: name,
		Path: path,
		Kind: kind,
	}
	switch kind {
	case "local":
		f.FS = os.DirFS(path)
	case "bundle":
		f.FS, err = archiver.FileSystem(context.Background(), path)
		if err != nil {
			return FS{}, err
		}
	}
	return f, nil
}
