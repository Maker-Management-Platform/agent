package models

import (
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/eduardooliveira/stLib/core/utils"
	"golang.org/x/exp/slices"
)

type ProjectAsset struct {
	SHA1         string        `json:"sha1" toml:"sha1" form:"sha1" query:"sha1" gorm:"primaryKey"`
	Name         string        `json:"name" toml:"name" form:"name" query:"name"`
	ProjectUUID  string        `json:"project_uuid" toml:"project_uuid" form:"project_uuid" query:"project_uuid" gorm:"primaryKey"`
	Project      *Project      `json:"-" toml:"-" form:"-" query:"-" gorm:"foreignKey:ProjectUUID"`
	Size         int64         `json:"size" toml:"size" form:"size" query:"size"`
	ModTime      time.Time     `json:"mod_time" toml:"mod_time" form:"mod_time" query:"mod_time"`
	AssetType    string        `json:"asset_type" toml:"asset_type" form:"asset_type" query:"asset_type"`
	Extension    string        `json:"extension" toml:"extension" form:"extension" query:"extension"`
	MimeType     string        `json:"mime_type" toml:"mime_type" form:"mime_type" query:"mime_type"`
	Model        *ProjectModel `json:"model" toml:"model" form:"model" query:"model"`
	ProjectImage *ProjectImage `json:"project_image" toml:"project_image" form:"project_image" query:"project_image"`
	ProjectFile  *ProjectFile  `json:"project_file" toml:"project_file" form:"project_file" query:"project_file"`
	Slice        *ProjectSlice `json:"slice" toml:"slice" form:"slice" query:"slice"`
}

func NewProjectAsset(fileName string, project *Project, file *os.File) (*ProjectAsset, []*ProjectAsset, error) {
	var asset = &ProjectAsset{
		Name:        fileName,
		ProjectUUID: project.UUID,
	}
	fullFilePath := utils.ToLibPath(fmt.Sprintf("%s/%s", project.FullPath(), fileName))

	var err error
	var nestedAssets []*ProjectAsset
	stat, err := file.Stat()
	if err != nil {
		return nil, nil, err
	}
	asset.Size = stat.Size()
	asset.ModTime = stat.ModTime()

	asset.Extension = filepath.Ext(fileName)
	asset.MimeType = mime.TypeByExtension(asset.Extension)
	asset.SHA1, err = utils.GetFileSha1(fullFilePath)
	if err != nil {
		return nil, nil, err
	}
	if slices.Contains(ModelExtensions, strings.ToLower(asset.Extension)) {
		asset.AssetType = ProjectModelType
		asset.Model, nestedAssets, err = NewProjectModel(fileName, asset, project, file)
	} else if slices.Contains(ImageExtensions, strings.ToLower(asset.Extension)) {
		asset.AssetType = ProjectImageType
		asset.ProjectImage, nestedAssets, err = NewProjectImage(fileName, asset, project, file)
	} else if slices.Contains(SliceExtensions, strings.ToLower(asset.Extension)) {
		asset.AssetType = ProjectSliceType
		asset.Slice, nestedAssets, err = NewProjectSlice(fileName, asset, project, file)
	} else {
		asset.AssetType = ProjectFileType
		asset.ProjectFile, nestedAssets, err = NewProjectFile(fileName, asset, project, file)
	}

	return asset, nestedAssets, err
}
