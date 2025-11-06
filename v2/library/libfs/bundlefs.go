package libfs

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"path/filepath"

	"github.com/eduardooliveira/stLib/v2/config"
	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/utils"
	"github.com/mholt/archiver/v4"
)

type bundleFS struct {
	lFS
	parentFS       LibFS
	bundleLocation string
}

func resolveBundleFS(ctx context.Context, asset entities.Asset) (LibFS, error) {
	if asset.FSKind == "bundle" {
		return resolveParentFS(ctx, asset)

	} else {
		return fileSystems[asset.FSName], nil
	}
}

func resolveParentFS(ctx context.Context, asset entities.Asset) (LibFS, error) {
	if asset.FSKind == "bundle" {
		var parentNode = utils.VoZ(asset.Parent)
		for parentNode.ID != asset.Root {
			parentNode = utils.VoZ(parentNode.Parent)
		}
		parent, err := resolveParentFS(ctx, parentNode)
		if err != nil {
			return nil, err
		}
		return newBundleFS(ctx, parent, parentNode)
	} else {
		return fileSystems[asset.FSName], nil
	}
}

func newBundleFS(ctx context.Context, parentFS LibFS, asset entities.Asset) (LibFS, error) {
	var path = utils.VoZ(asset.Path)
	var base = filepath.Base(path)
	switch filepath.Ext(base) {
	case ".zip", ".rar", ".7z", ".tar", ".3mf":
	default:
		return nil, errors.New("unsupported bundle type")
	}

	if parentFS.Kind() == "local" {
		bfs, err := archiver.FileSystem(ctx, filepath.Join(parentFS.GetLocation(), path))
		if err != nil {
			return nil, err
		}
		return bundleFS{
			lFS: lFS{
				FileSystem: config.FileSystem{
					Name: base,
					Path: asset.ID,
					Kind: "bundle",
				},
				FS:           bfs,
				discovarable: true,
			},
			parentFS:       parentFS,
			bundleLocation: filepath.Join(parentFS.GetLocation(), path),
		}, nil
	}

	tempFS, iErr := GetLibFS("temp")
	if iErr != nil {
		return nil, iErr
	}
	if _, err := fs.Stat(tempFS, base); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			writer, iErr := tempFS.Create(base)
			if iErr != nil {
				return nil, iErr
			}
			f, iErr := parentFS.Open(path)
			if iErr != nil {
				return nil, iErr
			}
			defer f.Close()
			_, iErr = io.Copy(writer, f)
			if iErr != nil {
				return nil, iErr
			}
		}
	}

	bfs, err := archiver.FileSystem(ctx, filepath.Join(tempFS.GetLocation(), base))

	if err != nil {
		return nil, err
	}
	return bundleFS{
		lFS: lFS{
			FileSystem: config.FileSystem{
				Name: parentFS.GetName(),
				Path: asset.ID,
				Kind: "bundle",
			},
			FS:           bfs,
			discovarable: true,
		},
		parentFS:       parentFS,
		bundleLocation: filepath.Join(tempFS.GetLocation(), base),
	}, nil
}

func (fs bundleFS) GetFS() fs.FS {
	return fs.lFS
}
func (fs bundleFS) GetName() string {
	return fs.Name
}
func (fs bundleFS) GetRoot() string {
	return fs.lFS.Path
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

func (fs bundleFS) Remove(name string) error {
	return errors.New("remove not supported")
}
