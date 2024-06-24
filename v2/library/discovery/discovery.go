package discovery

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/library/svc"
)

type discoverer struct {
	root         string
	l            *slog.Logger
	rootAsset    *entities.Asset
	currentAsset *entities.Asset
}

func New(root string) *discoverer {

	return &discoverer{
		root: root,
		l:    slog.With("module", "discovery").With("root", root),
	}
}

func (d *discoverer) ProcessPath(path string, parent *entities.Asset) (asset *entities.Asset, err error) {
	pathInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	rel, err := filepath.Rel(d.root, path)
	if err != nil {
		return nil, err
	}

	if parent == nil {
		d.rootAsset = asset
	}

	asset, _, err = svc.NewAssetFromRootPath(d.root, rel, pathInfo.IsDir(), parent)
	if err != nil {
		return nil, err
	}

	_, err = svc.SaveAsset(*asset)
	if err != nil {
		d.l.Error("Error saving asset", "error", err)
	}

	if pathInfo.IsDir() {
		files, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			a, err := d.ProcessPath(filepath.Join(path, file.Name()), asset)
			if err != nil {
				return nil, err
			}
			asset.NestedAssets = append(asset.NestedAssets, a)
		}
	}

	//d.l.Info("Asset saved", "asset", *aa.Path)

	return asset, nil
}

func (d *discoverer) Run() error {
	d.l.Info("Discovering assets")
	_, err := d.ProcessPath(d.root, nil)
	j, _ := json.MarshalIndent(d.rootAsset, "", "  ")
	os.WriteFile("assets.json", j, 0644)

	if err != nil {
		d.l.Error("Error discovering assets", "error", err)
	}

	qwe, err := svc.GetAssetByRootAndPath(d.root, ".", true)

	j, _ = json.MarshalIndent(qwe, "", "  ")
	os.WriteFile("assets2.json", j, 0644)

	return err
}
