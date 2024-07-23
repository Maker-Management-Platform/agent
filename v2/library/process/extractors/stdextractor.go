package extractors

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/eduardooliveira/stLib/v2/config"
	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/utils"
	"github.com/mholt/archiver/v4"
)

type StdExtractor struct {
}

func (t *StdExtractor) Extract(asset *entities.Asset) ([]*entities.Asset, error) {
	rtn := make([]*entities.Asset, 0)

	fsys, err := archiver.FileSystem(context.Background(), filepath.Join(*asset.Root, *asset.Path))
	if err != nil {
		return nil, err
	}
	imgRoot := filepath.Join(config.Cfg.Core.DataFolder, "img")
	parentDir := filepath.Join(imgRoot, *asset.ParentID)
	if err := utils.CreateFolder(parentDir); err != nil {
		return nil, err
	}

	var biggestImage int64
	var biggestImgFile string
	var files []string

	fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		files = append(files, path)
		ext := filepath.Ext(d.Name())
		if config.Cfg.Library.AssetTypes.ByExtension(ext).Name == "image" {
			info, err := d.Info()
			if err != nil {
				return err
			}

			if info.Size() > biggestImage {
				biggestImage = info.Size()
				biggestImgFile = path
			}
		}

		return nil
	})

	for _, file := range files {
		var za *entities.Asset
		fName := filepath.Base(file)
		ext := filepath.Ext(fName)
		name := fmt.Sprintf("%s.e%s", strings.TrimSuffix(fName, ext), ext)
		if file == biggestImgFile {
			f, err := fsys.Open(file)
			if err != nil {
				return nil, err
			}
			defer f.Close()
			if err := utils.SaveFile(filepath.Join(parentDir, name), f); err != nil {
				return nil, err
			}
			za = entities.NewAssetFromRootPath(imgRoot, filepath.Join(*asset.ParentID, name), false, asset)
			if asset.Thumbnail == nil {
				asset.Thumbnail = &za.ID
			}
		} else {
			za = entities.NewBundledAsset(asset, file)
		}

		rtn = append(rtn, za)
	}

	asset.NodeKind = utils.Ptr("bundle")
	return rtn, nil
}

// TODO: propagate context
func (t *StdExtractor) ExtractBundled(asset *entities.Asset) error {
	if asset.Parent == nil || asset.Parent.Root == nil || asset.Parent.Path == nil {
		return errors.New("invalid parent asset")
	}
	fsys, err := archiver.FileSystem(context.Background(), filepath.Join(*asset.Root, *asset.Parent.Path))
	if err != nil {
		return err
	}

	f, err := fsys.Open(utils.VoZ(asset.Path))
	if err != nil {
		return err
	}
	defer f.Close()
	parentDir := filepath.Dir(*asset.Parent.Path)
	assetBase := filepath.Base(*asset.Path)

	if err := utils.SaveFile(filepath.Join(*asset.Parent.Root, parentDir, assetBase), f); err != nil {
		return err
	}

	asset.Path = utils.Ptr(filepath.Join(parentDir, assetBase))
	asset.Root = asset.Parent.Root
	asset.NodeKind = utils.Ptr(entities.NodeKindFile)

	kind := config.Cfg.Library.AssetTypes.ByExtension(*asset.Extension).Name
	asset.Kind = utils.Ptr(kind)
	if kind == "image" {
		asset.Thumbnail = utils.Ptr(asset.ID)
	}

	return nil
}
