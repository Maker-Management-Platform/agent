package discovery

import (
	"io/fs"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/eduardooliveira/stLib/v2/config"
	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/library/libfs"
	"github.com/eduardooliveira/stLib/v2/library/process"
	"github.com/eduardooliveira/stLib/v2/library/repo"
	"github.com/eduardooliveira/stLib/v2/utils"
)

type discProc struct {
	discoverer *Discoverer
	l          *slog.Logger
	fs         libfs.LibFS
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

func (d Discoverer) DiscoverFS(fs libfs.LibFS) *discProc {
	dp := &discProc{
		discoverer: &d,
		l:          d.l.With("fs", fs.GetName()),
		fs:         fs,
	}

	return dp
}

func (d *discProc) processPath(currFS libfs.LibFS, path string, parent *entities.Asset) (asset *entities.Asset, err error) {

	if path != "." && d.shouldSkipFile(path) {
		return nil, nil
	}

	pathInfo, err := fs.Stat(currFS, path)
	if err != nil {
		return nil, err
	}

	asset = entities.NewAsset(currFS, path, pathInfo.IsDir(), parent)

	asset.SeenOnScan = utils.Ptr(true)

	err = d.discoverer.repo.SaveAsset(*asset)
	if err != nil {
		d.l.Error("Error saving asset", "error", err)
	}

	if pathInfo.IsDir() || libfs.IsBundle(path) {

		innerFS := currFS
		var files []fs.DirEntry
		if libfs.IsBundle(path) {
			innerFS, err = libfs.GetBundleFS(currFS, path)
			if err != nil {
				return nil, err
			}
			files, err = fs.ReadDir(innerFS, ".")
			if err != nil {
				return nil, err
			}
			path = "."
		} else {
			files, err = fs.ReadDir(currFS, path)
			if err != nil {
				return nil, err
			}
		}

		for _, file := range files {
			_, err := d.processPath(innerFS, filepath.Join(path, file.Name()), asset)
			if err != nil {
				return nil, err
			}
		}
	}
	if !pathInfo.IsDir() || (!libfs.IsBundle(path) || config.Cfg.Library.RenderBundles) {
		d.discoverer.processor.Process(asset)
	}

	return asset, nil
}

func (d *discProc) Run() error {
	d.l.Info("Discovering assets")

	err := d.discoverer.repo.SetDirtyFS(d.fs.GetName())
	if err != nil {
		return err
	}

	_, err = d.processPath(d.fs, ".", nil)

	if err != nil {
		d.l.Error("Error discovering assets", "error", err)
	}

	err = d.discoverer.repo.DeleteUnSeenInFS(d.fs.GetName())
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
