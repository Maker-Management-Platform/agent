package discovery

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/library/process"
	"github.com/eduardooliveira/stLib/v2/library/repo"
)

type discProc struct {
	d         *Discoverer
	l         *slog.Logger
	root      string
	rootAsset *entities.Asset
}

type Discoverer struct {
	l *slog.Logger
	p *process.Processor
	r *repo.AssetRepo
}

func New(p *process.Processor, r *repo.AssetRepo) *Discoverer {
	return &Discoverer{
		p: p,
		r: r,
		l: slog.With("module", "discovery"),
	}
}

func (d *Discoverer) Get(root string) *discProc {
	return &discProc{
		d:    d,
		l:    d.l.With("root", root),
		root: root,
	}
}

func (d *discProc) ProcessPath(path string, parent *entities.Asset) (asset *entities.Asset, err error) {
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

	asset = entities.NewAssetFromRootPath(d.root, rel, pathInfo.IsDir(), parent)

	err = d.d.r.SaveAsset(*asset)
	if err != nil {
		d.l.Error("Error saving asset", "error", err)
	}

	d.d.p.Process(asset)

	if pathInfo.IsDir() {
		files, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			_, err := d.ProcessPath(filepath.Join(path, file.Name()), asset)
			if err != nil {
				return nil, err
			}
		}
	}

	return asset, nil
}

func (d *discProc) Run() error {
	d.l.Info("Discovering assets")
	_, err := d.ProcessPath(d.root, nil)

	if err != nil {
		d.l.Error("Error discovering assets", "error", err)
	}

	return err
}
