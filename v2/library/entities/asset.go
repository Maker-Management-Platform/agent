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
	ID           string     `form:"id" gorm:"primaryKey"`
	Label        *string    `form:"label"`
	Description  *string    `form:"description"`
	Path         *string    `form:"path"`
	Root         *string    `form:"root"`
	Extension    *string    `form:"extension"`
	Kind         *string    `form:"kind"`
	NodeKind     *string    `form:"nodeKind"`
	ParentID     *string    `form:"parentID"`
	Parent       *Asset     `form:"-"`
	NestedAssets []*Asset   `form:"nestedAssets" gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE;"`
	Thumbnail    *string    `form:"thumbnail"`
	SeenOnScan   *bool      `form:"seenOnScan"`
	Properties   Properties `form:"properties"`
	Tags         []*Tag     `form:"tags" gorm:"many2many:asset_tags"`
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
		asset.Kind = utils.Ptr("dir")
		return asset
	}

	asset.NodeKind = utils.Ptr(NodeKindFile)
	kind := config.Cfg.Library.AssetTypes.ByExtension(*asset.Extension)
	asset.Kind = utils.Ptr(kind.Name)

	if *asset.Kind == "image" {
		asset.Thumbnail = utils.Ptr(asset.ID)
	}

	return asset
}

func (a *Asset) bubbleThumbnail(tx *gorm.DB) error {
	if a.Thumbnail == nil || a.ParentID == nil {
		slog.Debug("no thumbnail or parent", "asset", utils.VoZ(a.Path), "thumbnail", utils.VoZ(a.Thumbnail), "parent", utils.VoZ(a.ParentID))
		return nil
	}

	var parent = &Asset{}

	if q := tx.Model(Asset{}).Where(Asset{ID: *a.ParentID}).First(parent); q.Error != nil {
		slog.With("asset", *a.Path).With("context", "bubbleThumbnail").With("error", q.Error).Debug("error loading asset")
		return q.Error
	}

	if parent.Thumbnail == nil {
		parent.Thumbnail = a.Thumbnail
		slog.Debug("bubbling thumbnail", "asset", utils.VoZ(a.Path), "thumbnail", utils.VoZ(a.Thumbnail), "parent", utils.VoZ(a.ParentID))
		if q := tx.Save(parent); q.Error != nil {
			slog.With("asset", *parent.Path).With("context", "bubbleThumbnail").With("error", q.Error).Error("error saving asset")
			return q.Error
		}
	}

	return nil
}

func (a *Asset) AfterSave(tx *gorm.DB) error {
	if a.ID == "" {
		return nil
	}
	slog.Debug("AfterSave", "asset", utils.VoZ(a.Path))
	return a.bubbleThumbnail(tx)
}
