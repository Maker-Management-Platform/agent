package entities

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/fs"
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	"github.com/eduardooliveira/stLib/v2/config"
	"github.com/eduardooliveira/stLib/v2/utils"
	"gorm.io/gorm"
)

type LibFS interface {
	GetFS() fs.FS
	GetName() string
	GetLocation() string
	GetRoot() string
	Kind() string
	Open(name string) (fs.File, error)
	Writable() bool
	Create(name string) (io.WriteCloser, error)
	Mkdir(name string) error
	Remove(name string) error
	IsBundle(path string) bool
}

type NodeKind string

const (
	NodeKindRoot    NodeKind = "root"
	NodeKindFile    NodeKind = "file"
	NodeKindDir     NodeKind = "dir"
	NodeKindBundle  NodeKind = "bundle"
	NodeKindBundled NodeKind = "bundled"
)

type Asset struct {
	ID           string     `query:"id" in:"form=id" gorm:"primaryKey"`
	Label        *string    `query:"label" in:"form=label"`
	Description  *string    `query:"description" in:"form=description"`
	Path         *string    `query:"path" in:"form=path"`
	Root         string     `query:"root" in:"form=root"`
	FSKind       string     `query:"fsKind" in:"form=fsKind"`
	FSName       string     `query:"fsName" in:"form=fsName"`
	Extension    *string    `query:"extension" in:"form=extension"`
	Kind         *string    `query:"kind" in:"form=kind"`
	NodeKind     NodeKind   `query:"nodeKind" in:"form=nodeKind"`
	ParentID     *string    `query:"parentID" in:"form=parentID"`
	Parent       *Asset     `in:"form=-"`
	NestedAssets []*Asset   `query:"nestedAssets" in:"form=nestedAssets" gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE;"`
	Thumbnail    *string    `query:"thumbnail" in:"form=thumbnail"`
	SeenOnScan   *bool      `query:"seenOnScan" in:"form=seenOnScan"`
	Properties   Properties `query:"properties" in:"form=properties"`
	Tags         []*Tag     `query:"tags" in:"form=tags" gorm:"many2many:asset_tags"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewAsset(fs LibFS, path string, isDir bool, parent *Asset) *Asset {
	ext := filepath.Ext(path)

	data := []byte(filepath.Join(fs.GetName(), fs.GetRoot(), path))
	md5Hash := md5.Sum(data)

	var asset = &Asset{
		ID:        hex.EncodeToString(md5Hash[:]),
		Path:      utils.Ptr(path),
		Root:      fs.GetRoot(),
		FSName:    fs.GetName(),
		FSKind:    fs.Kind(),
		Label:     utils.Ptr(strings.TrimSuffix(filepath.Base(path), ext)),
		Extension: utils.Ptr(ext),
	}
	if parent != nil {
		asset.Parent = parent
		asset.ParentID = &parent.ID
	}

	if fs.IsBundle(path) {
		asset.Kind = utils.Ptr("bundle")
		asset.NodeKind = NodeKindBundle
	} else if isDir {
		if parent == nil {
			asset.NodeKind = NodeKindRoot
			asset.Label = utils.Ptr(fs.GetName())
		} else {
			asset.NodeKind = NodeKindDir
		}
	} else {
		asset.NodeKind = NodeKindFile
	}
	if fs.Kind() == "bundle" {
		asset.NodeKind = NodeKindBundled
	}

	if isDir {
		asset.Kind = utils.Ptr("dir")
		return asset
	}

	if asset.Kind == nil {
		kind := config.Cfg.Library.AssetTypes.ByExtension(*asset.Extension)
		asset.Kind = utils.Ptr(kind.Name)
	}
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
