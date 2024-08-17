package library

import (
	"log/slog"
	"net/http"

	"github.com/eduardooliveira/stLib/v2/config"
	"github.com/eduardooliveira/stLib/v2/library/discovery"
	"github.com/eduardooliveira/stLib/v2/library/process"
	"github.com/eduardooliveira/stLib/v2/library/repo"
	"github.com/eduardooliveira/stLib/v2/library/web"
	"golang.org/x/sync/errgroup"
)

type Library struct {
	r *repo.AssetRepo
	p *process.Processor
	d *discovery.Discoverer
}

func New() (*Library, http.Handler, http.Handler, error) {
	lib := &Library{}

	var err error

	lib.r, err = repo.New()
	if err != nil {
		return nil, nil, nil, err
	}

	lib.p, err = process.New(lib.r)
	if err != nil {
		return nil, nil, nil, err
	}

	lib.d = discovery.New(lib.p, lib.r)

	webH, err := web.New(lib.r, lib.p)
	if err != nil {
		return nil, nil, nil, err
	}

	/*apiH, err := api.NewChi(lib.r)
	if err != nil {
		return nil, nil, nil, err
	}*/

	return lib, webH, nil, nil
}
func (l Library) Scan() error {
	eg := errgroup.Group{}
	if len(config.Cfg.Library.Paths) == 0 {
		slog.Warn("No library paths configured")
	}
	for _, path := range config.Cfg.Library.Paths {
		eg.Go(l.d.Get(path).Run)
	}
	return eg.Wait()
}

func (l Library) ScanFS() error {
	eg := errgroup.Group{}
	if len(config.Cfg.Library.FileSystems) == 0 || (len(config.Cfg.Library.FileSystems) == 1 && config.Cfg.Library.FileSystems[0].Path == "change_me") {
		slog.Warn("invalid library file systems configured")
		return nil
	}
	for _, cfs := range config.Cfg.Library.FileSystems {
		eg.Go(l.d.GetForFS(cfs).Run)
	}
	return eg.Wait()
}

func (l *Library) ScanAsync() {
	go func() {
		if err := l.ScanFS(); err != nil {
			slog.Error("Error scanning library", "error", err)
		}
	}()
}
