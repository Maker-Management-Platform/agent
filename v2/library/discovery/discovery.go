package discovery

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/eduardooliveira/stLib/v2/config"
	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/library/process"
	"github.com/eduardooliveira/stLib/v2/library/repo"
	"github.com/eduardooliveira/stLib/v2/utils"
)

type discProc struct {
	discoverer *Discoverer
	l          *slog.Logger
	root       string
	rootAsset  *entities.Asset
}

type Discoverer struct {
	l         *slog.Logger
	processor *process.Processor
	repo      *repo.AssetRepo
}

func New(p *process.Processor, r *repo.AssetRepo) *Discoverer {
	return &Discoverer{
		processor: p,
		repo:      r,
		l:         slog.With("module", "discovery"),
	}
}

func (d *Discoverer) Get(root string) *discProc {
	return &discProc{
		discoverer: d,
		l:          d.l.With("root", root),
		root:       root,
	}
}

func (d *discProc) ProcessPath(path string, parent *entities.Asset) (asset *entities.Asset, err error) {

	if d.shouldSkipFile(path) {
		return nil, nil
	}

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

	asset.SeenOnScan = utils.Ptr(true)

	err = d.discoverer.repo.SaveAsset(*asset)
	if err != nil {
		d.l.Error("Error saving asset", "error", err)
	}

	d.discoverer.processor.Process(asset)

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

	err := d.discoverer.repo.SetDirtyRoot(d.root)
	if err != nil {
		return err
	}

	_, err = d.ProcessPath(d.root, nil)

	if err != nil {
		d.l.Error("Error discovering assets", "error", err)
	}

	err = d.discoverer.repo.DeleteUnSeenInRoot(d.root)
	if err != nil {
		return err
	}

	return err
}

func (d *discProc) shouldSkipFile(name string) bool {

	if strings.HasPrefix(name, ".") {
		if config.Cfg.Library.IgnoreDotFiles {
			return true
		}
	}

	for _, blacklist := range config.Cfg.Library.Blacklist {
		if strings.HasSuffix(name, blacklist) {
			return true
		}
	}

	return false
}
