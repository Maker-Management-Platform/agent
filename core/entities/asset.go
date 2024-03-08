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
	ID          string          `json:"id" toml:"id" form:"id" query:"id" gorm:"primaryKey"`
	Name        string          `json:"name" toml:"name" form:"name" query:"name"`
	Label       string          `json:"label" toml:"label" form:"label" query:"label"`
	Origin      string          `json:"origin" toml:"origin" form:"origin" query:"origin"`
	ProjectUUID string          `json:"project_uuid" toml:"project_uuid" form:"project_uuid" query:"project_uuid"`
	project     *Project        `json:"-" toml:"-" form:"-" query:"-" gorm:"foreignKey:ProjectUUID"`
	Size        int64           `json:"size" toml:"size" form:"size" query:"size"`
	ModTime     time.Time       `json:"mod_time" toml:"mod_time" form:"mod_time" query:"mod_time"`
	AssetType   string          `json:"asset_type" toml:"asset_type" form:"asset_type" query:"asset_type"`
	Extension   string          `json:"extension" toml:"extension" form:"extension" query:"extension"`
	MimeType    string          `json:"mime_type" toml:"mime_type" form:"mime_type" query:"mime_type"`
	ImageID     string          `json:"image_id" toml:"image_id" form:"image_id" query:"image_id"`
	Properties  AssetProperties `json:"properties" toml:"properties" form:"properties" query:"properties"`
}

func NewProjectAsset2(fileName string, label string, project *Project, origin string) (*ProjectAsset, error) {
	var asset = &ProjectAsset{
		Name:        fileName,
		Label:       label,
		Origin:      origin,
		ProjectUUID: project.UUID,
		project:     project,
		Properties:  make(map[string]any),
	}

	var fullFilePath string
	if origin == "fs" {
		fullFilePath = utils.ToLibPath(path.Join(project.FullPath(), fileName))
	} else {
		fullFilePath = utils.ToAssetsPath(project.UUID, fileName)
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
	asset.ID, err = assetSha(project.UUID, fileName, fullFilePath)
	if err != nil {
		return nil, err
	}
	return asset, err
}

func assetSha(projectUuid string, assetName string, fullFilePath string) (string, error) {
	fSha512, err := utils.GetFileSha512(fullFilePath)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", sha1.Sum([]byte(fmt.Sprintf("%s%s%s", projectUuid, assetName, fSha512)))), nil
}
