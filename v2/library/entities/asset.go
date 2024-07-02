package entities

import (
	"crypto/md5"
	"encoding/hex"
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	"github.com/eduardooliveira/stLib/v2/config"
	"github.com/eduardooliveira/stLib/v2/utils"
	"gorm.io/gorm"
)

const (
	NodeKindRoot   = "root"
	NodeKindFile   = "file"
	NodeKindDir    = "dir"
	NodeKindBundle = "bundle"
)

type Asset struct {
	ID           string     `json:"id" gorm:"primaryKey"`
	Label        *string    `json:"label"`
	Path         *string    `json:"path"`
	Root         *string    `json:"root"`
	Extension    *string    `json:"extension"`
	Kind         *string    `json:"kind"`
	NodeKind     *string    `json:"nodeKind"`
	ParentID     *string    `json:"parentID"`
	Parent       *Asset     `json:"-"`
	NestedAssets []*Asset   `json:"nestedAssets" gorm:"foreignKey:ParentID"`
	Thumbnail    *string    `json:"thumbnail"`
	SeenOnScan   *bool      `json:"seenOnScan"`
	Properties   Properties `json:"properties"`
	Tags         []*Tag     `json:"tags" gorm:"many2many:asset_tags"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewAssetFromRootPath(root, path string, isDir bool, parent *Asset) *Asset {
	ext := filepath.Ext(path)

	data := []byte(filepath.Join(root, path)) // Convert the path string to bytes
	md5Hash := md5.Sum(data)                  // Generate the MD5 hash

	var asset = &Asset{
		ID:        hex.EncodeToString(md5Hash[:]),
		Root:      utils.Ptr(root),
		Path:      utils.Ptr(path),
		Label:     utils.Ptr(strings.TrimSuffix(filepath.Base(path), ext)),
		Extension: utils.Ptr(ext),
	}
	if parent != nil {
		asset.Parent = parent
		asset.ParentID = &parent.ID
	}
	if isDir {
		if parent == nil {
			asset.NodeKind = utils.Ptr(NodeKindRoot)
		} else {
			asset.NodeKind = utils.Ptr(NodeKindDir)
		}
		return asset
	}
	kind := config.Cfg.Library.AssetTypes.ByExtension(*asset.Extension)
	asset.Kind = utils.Ptr(kind.Name)

	if *asset.Kind == "image" {
		asset.Thumbnail = utils.Ptr(asset.ID)
	}

	return asset
}

func (a *Asset) bubbleThumbnail(tx *gorm.DB) error {
	if a.Thumbnail == nil || a.ParentID == nil {
		slog.Debug("no thumbnail or parent", "asset", *a.Path, "thumbnail", a.Thumbnail, "parent", a.ParentID)
		return nil
	}

	var parent = &Asset{}

	if q := tx.Model(Asset{}).Where(Asset{ID: *a.ParentID}).First(parent); q.Error != nil {
		slog.With("asset", *a.Path).With("context", "bubbleThumbnail").With("error", q.Error).Debug("error loading asset")
		return q.Error
	}

	if parent.Thumbnail == nil {
		parent.Thumbnail = a.Thumbnail
		slog.Debug("bubbling thumbnail", "asset", *a.Path, "thumbnail", *a.Thumbnail, "parent", *parent.Path)
		if q := tx.Save(parent); q.Error != nil {
			slog.With("asset", *parent.Path).With("context", "bubbleThumbnail").With("error", q.Error).Error("error saving asset")
			return q.Error
		}
	}

	return nil
}

func (a *Asset) AfterSave(tx *gorm.DB) error {
	slog.Debug("AfterSave", "asset", *a.Path)
	return a.bubbleThumbnail(tx)
}
