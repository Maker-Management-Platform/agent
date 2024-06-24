package entities

import (
	"log/slog"
	"time"

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
	Name         *string    `json:"name"`
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

func NewAsset(root, path string) *Asset {
	return &Asset{
		Root: &root,
		Path: &path,
	}
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
