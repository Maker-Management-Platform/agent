package entities

import (
	"crypto/sha1"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/eduardooliveira/stLib/core/utils"
	"golang.org/x/exp/slices"
)

type AssetProperties map[string]any

func (n *AssetProperties) Scan(src interface{}) error {
	str, ok := src.(string)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSON string:", src))
	}
	a := json.Unmarshal([]byte(str), &n)
	return a
}

func (n AssetProperties) Value() (driver.Value, error) {
	val, err := json.Marshal(n)
	return string(val), err
}

type ProjectAsset struct {
	ID           string          `json:"id" toml:"id" form:"id" query:"id" gorm:"primaryKey"`
	Name         string          `json:"name" toml:"name" form:"name" query:"name"`
	Generated    bool            `json:"generated" toml:"generated" form:"generated" query:"generated"`
	ProjectUUID  string          `json:"project_uuid" toml:"project_uuid" form:"project_uuid" query:"project_uuid"`
	project      *Project        `json:"-" toml:"-" form:"-" query:"-" gorm:"foreignKey:ProjectUUID"`
	Size         int64           `json:"size" toml:"size" form:"size" query:"size"`
	ModTime      time.Time       `json:"mod_time" toml:"mod_time" form:"mod_time" query:"mod_time"`
	AssetType    string          `json:"asset_type" toml:"asset_type" form:"asset_type" query:"asset_type"`
	Extension    string          `json:"extension" toml:"extension" form:"extension" query:"extension"`
	MimeType     string          `json:"mime_type" toml:"mime_type" form:"mime_type" query:"mime_type"`
	ImageID      string          `json:"image_id" toml:"image_id" form:"image_id" query:"image_id"`
	Properties   AssetProperties `json:"properties" toml:"properties" form:"properties" query:"properties"`
	State        string          `json:"state" toml:"state" form:"state" query:"state"`
	Model        *ProjectModel   `json:"model" toml:"model" form:"model" query:"model"`
	ProjectImage *ProjectImage   `json:"project_image" toml:"project_image" form:"project_image" query:"project_image"`
	ProjectFile  *ProjectFile    `json:"project_file" toml:"project_file" form:"project_file" query:"project_file"`
	Slice        *ProjectSlice   `json:"slice" toml:"slice" form:"slice" query:"slice"`
}

var GeneratedExtensions = []string{".thumb.png", ".render.png"}

func NewProjectAsset(fileName string, project *Project, file *os.File) (*ProjectAsset, []*ProjectAsset, error) {
	var asset = &ProjectAsset{
		Name:        fileName,
		Generated:   false,
		ProjectUUID: project.UUID,
		project:     project,
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

	asset.Extension = strings.ToLower(filepath.Ext(fileName))
	asset.MimeType = mime.TypeByExtension(asset.Extension)
	asset.ID, err = assetSha1(project.UUID, fileName, fullFilePath)
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
	for _, ext := range GeneratedExtensions {
		if strings.HasSuffix(asset.Name, ext) {
			asset.Generated = true
		}
		for _, nestedAsset := range nestedAssets {
			if strings.HasSuffix(nestedAsset.Name, ext) {
				nestedAsset.Generated = true
			}
		}
	}

	return asset, nestedAssets, err
}
func NewProjectAsset2(fileName string, project *Project, generated bool) (*ProjectAsset, error) {
	var asset = &ProjectAsset{
		Name:        fileName,
		Generated:   generated,
		ProjectUUID: project.UUID,
		project:     project,
		Properties:  make(map[string]any),
	}

	var fullFilePath string
	if generated {
		fullFilePath = utils.ToGeneratedPath(fileName)
	} else {
		fullFilePath = utils.ToLibPath(path.Join(project.FullPath(), fileName))
	}

	var err error
	file, err := os.Open(fullFilePath)
	if err != nil {
		log.Println("failed to open file", err)
		return nil, err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	asset.Size = stat.Size()
	asset.ModTime = stat.ModTime()

	asset.Extension = strings.ToLower(filepath.Ext(fileName))
	asset.MimeType = mime.TypeByExtension(asset.Extension)
	asset.ID, err = assetSha1(project.UUID, fileName, fullFilePath)
	if err != nil {
		return nil, err
	}
	return asset, err
}

func assetSha1(projectUuid string, assetName string, fullFilePath string) (string, error) {
	fSha512, err := utils.GetFileSha512(fullFilePath)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", sha1.Sum([]byte(fmt.Sprintf("%s%s%s", projectUuid, assetName, fSha512)))), nil
}
