package extractors

import (
	"archive/zip"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/eduardooliveira/stLib/v2/config"
	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/utils"
)

type ThreeMFExtractor struct {
}

func (t *ThreeMFExtractor) Extract(asset *entities.Asset) ([]*entities.Asset, error) {
	rtn := make([]*entities.Asset, 0)

	archive, err := zip.OpenReader(filepath.Join(*asset.Root, *asset.Path))
	if err != nil {
		return nil, err
	}
	defer archive.Close()

	imgRoot := filepath.Join(config.Cfg.Core.DataFolder, "img")
	parentDir := filepath.Join(imgRoot, *asset.ParentID)
	if err := utils.CreateFolder(parentDir); err != nil {
		return nil, err
	}

	thumbnail := t.filterThumb(archive.File)

	// elect a tumbnail
	// itersate over content map assets and ectract tumb
	for _, f := range archive.File {

		if f.FileInfo().IsDir() {
			continue
		}

		if strings.Contains(f.Name, ".thumbnails/") {
			continue
		}

		fName := filepath.Base(f.Name)
		ext := filepath.Ext(fName)
		name := fmt.Sprintf("%s.e%s", strings.TrimSuffix(fName, ext), ext)

		var za *entities.Asset

		if f.Name == thumbnail {
			zipped, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer zipped.Close()
			if err := utils.SaveFile(filepath.Join(parentDir, name), zipped); err != nil {
				return nil, err
			}

			za = entities.NewAssetFromRootPath(imgRoot, filepath.Join(*asset.ParentID, name), false, asset)
			if asset.Thumbnail == nil {
				asset.Thumbnail = &za.ID
			}
		} else {
			za = entities.NewBundledAsset(asset, f.Name)
		}

		rtn = append(rtn, za)

	}

	asset.NodeKind = utils.Ptr("bundle")
	return rtn, nil
}

func (t *ThreeMFExtractor) filterThumb(files []*zip.File) string {
	var biggestSize uint64
	var biggestImage string

	for _, f := range files {
		// Ignore thumbnail since we should have the original image already
		if strings.Contains(f.Name, ".thumbnails/") {
			continue
		}

		ext := filepath.Ext(f.Name)
		if config.Cfg.Library.AssetTypes.ByExtension(ext).Name != "image" {
			continue
		}

		// try set makerworld thumbnail
		if strings.HasSuffix(f.Name, "Preview.webp") {
			return f.Name
		}

		if f.UncompressedSize64 > biggestSize {
			biggestSize = f.UncompressedSize64
			biggestImage = f.Name
		}
	}

	return biggestImage
}
