package svc

import (
	"errors"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/eduardooliveira/stLib/v2/config"
	"github.com/eduardooliveira/stLib/v2/database"
	"github.com/eduardooliveira/stLib/v2/library/enrichers"
	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/library/extractors"
	"github.com/eduardooliveira/stLib/v2/library/renderers"
	"github.com/eduardooliveira/stLib/v2/library/svc/data"
	"github.com/eduardooliveira/stLib/v2/utils"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type service struct {
	l *slog.Logger
}

var svc *service

func Init() error {
	svc = &service{
		l: slog.With("module", "svc"),
	}

	if err := database.DB.AutoMigrate(&entities.Asset{}); err != nil {
		svc.l.Error("failed to auto migrate asset", "error", err)
		return err
	}

	return nil
}

func NewAssetFromRootPath(root, path string, isDir bool, parent *entities.Asset) (asset *entities.Asset, found bool, err error) {
	asset = &entities.Asset{
		ID:   uuid.New().String(),
		Root: &root,
		Path: &path,
	}

	if a, err := GetAssetByRootAndPath(root, path, false); err == nil {
		asset = a
		found = true
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, err
	}

	if parent == nil {
		asset.NodeKind = utils.Ptr(entities.NodeKindRoot)
		return
	} else {
		asset.ParentID = utils.Ptr(parent.ID)
		parent.Parent = parent
	}

	if isDir {
		asset.NodeKind = utils.Ptr(entities.NodeKindDir)
		return
	}

	asset.Extension = utils.Ptr(strings.ToLower(filepath.Ext(path)))
	kind := config.Cfg.Library.AssetTypes.ByExtension(*asset.Extension)
	asset.Kind = utils.Ptr(kind.Name)

	if *asset.Kind == "image" {
		asset.Thumbnail = utils.Ptr(asset.ID)
	}

	eg := errgroup.Group{}

	if extractors.IsExtractable(asset) {
		asset.NodeKind = utils.Ptr(entities.NodeKindBundle)
	} else {
		asset.NodeKind = utils.Ptr(entities.NodeKindFile)
	}

	if enricher, ok := enrichers.Get(asset); ok {
		eg.Go(enricher.Enrich(asset))
	}
	if renderers.IsRenderable(asset) {
		eg.Go(renderers.GetRenderer(asset).Render(*asset, onRendered))
	}

	go func() {
		if err := eg.Wait(); err != nil {
			svc.l.With("asset", asset.Path).Error("Error processing asset", "error", err)
		}
	}()

	return
}

func onRendered(asset *entities.Asset, imgRoot, imgName string) error {
	rendered, _, err := NewAssetFromRootPath(imgRoot, imgName, false, asset)
	if err != nil {
		svc.l.With("asset", asset.Path).Error("Error getting rendered asset", "error", err)
		return err
	}
	if _, err := SaveAsset(*rendered); err != nil {
		svc.l.With("asset", asset.Path).Error("Error saving asset", "error", err)
		return err
	}

	return nil
}

func SaveAsset(a entities.Asset) (*entities.Asset, error) {
	if err := data.SaveAsset(a); err != nil {
		return nil, err
	}
	return data.GetAsset(a.ID, false)
}

func GetAssetByRootAndPath(root, path string, deep bool) (*entities.Asset, error) {
	return data.GetAssetByRootAndPath(root, path, deep)
}

func GetAsset(id string, deep bool) (*entities.Asset, error) {
	return data.GetAsset(id, deep)
}
