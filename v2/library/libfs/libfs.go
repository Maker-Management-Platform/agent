package libfs

import (
	"errors"
	"io"
	"io/fs"
	"path/filepath"
	"slices"

	"github.com/eduardooliveira/stLib/v2/config"
	"github.com/eduardooliveira/stLib/v2/utils"
)

type FS struct {
	fs.FS
	config.FileSystem
	discovarable bool
}

func (fs FS) isDiscovarable() bool {
	return fs.discovarable
}

type LibFS interface {
	GetFS() fs.FS
	GetName() string
	GetLocation() string
	Kind() string
	Open(name string) (fs.File, error)
	Writable() bool
	Create(name string) (io.WriteCloser, error)
	Mkdir(name string) error
	isDiscovarable() bool
	setDiscovarable(bool)
}

var (
	fileSystems   = make(map[string]LibFS)
	defaultFSName = ""
	bundleFSs     = []string{".zip", ".rar", ".7z", ".tar", ".3mf"}
)

func LoadFSs() error {
	if len(config.Cfg.Library.FileSystems) == 0 {
		return errors.New("no file systems found")
	}
	for _, cfgFS := range config.Cfg.Library.FileSystems {
		switch cfgFS.Kind {
		case "local":
			fileSystems[cfgFS.Name] = newLocalFS(cfgFS)
		case "gitfs":
			gitfs, err := newGitFS(cfgFS)
			if err != nil {
				return err
			}
			fileSystems[cfgFS.Name] = gitfs
		}
		if defaultFSName == "" || cfgFS.Default {
			defaultFSName = cfgFS.Name
		}
	}
	if fileSystems["cache"] == nil {
		if err := utils.CreateFolder(filepath.Join(config.Cfg.Core.DataFolder, "cache")); err != nil {
			return err
		}
		fileSystems["cache"] = newLocalFS(config.FileSystem{
			Name: "cache",
			Kind: "local",
			Path: filepath.Join(config.Cfg.Core.DataFolder, "cache"),
		})
	}
	if fileSystems["generated"] == nil {
		if err := utils.CreateFolder(filepath.Join(config.Cfg.Core.DataFolder, "generated")); err != nil {
			return err
		}
		fileSystems["generated"] = newLocalFS(config.FileSystem{
			Name: "generated",
			Kind: "local",
			Path: filepath.Join(config.Cfg.Core.DataFolder, "generated"),
		})
	}
	fileSystems["cache"].setDiscovarable(false)
	fileSystems["generated"].setDiscovarable(false)
	return nil
}

func GetDefaultFS() LibFS {
	return fileSystems[defaultFSName]
}

func GetFS(kind, name, path string) (LibFS, error) {
	if kind == "bundle" {
		return newBundleFS(fileSystems[name], path)
	}
	if _, ok := fileSystems[name]; !ok {
		return nil, errors.New("file system not found")
	}
	return fileSystems[name], nil
}

func GetLibFS(name string) (LibFS, error) {
	if _, ok := fileSystems[name]; !ok {
		return nil, errors.New("file system not found")
	}
	return fileSystems[name], nil
}

func GetFSs() map[string]LibFS {
	rtn := make(map[string]LibFS)
	for k, v := range fileSystems {
		if v.isDiscovarable() {
			rtn[k] = v
		}
	}

	return rtn
}

func IsBundle(path string) bool {
	return slices.Contains(bundleFSs, filepath.Ext(path))
}

func GetBundleFS(parentFS LibFS, path string) (LibFS, error) {
	return newBundleFS(parentFS, path)
}
