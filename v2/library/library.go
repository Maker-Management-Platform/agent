package library

import (
	"log/slog"

	"github.com/eduardooliveira/stLib/v2/config"
	"github.com/eduardooliveira/stLib/v2/library/api"
	"github.com/eduardooliveira/stLib/v2/library/discovery"
	"github.com/eduardooliveira/stLib/v2/library/process"
	"github.com/eduardooliveira/stLib/v2/library/repo"
	"github.com/eduardooliveira/stLib/v2/library/web"
	"github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"
)

type Library struct {
	r *repo.AssetRepo
	p *process.Processor
	d *discovery.Discoverer
}

func New(e *echo.Group) (*Library, error) {
	lib := &Library{}

	var err error

	lib.r, err = repo.New()
	if err != nil {
		return nil, err
	}

	lib.p, err = process.New(lib.r)
	if err != nil {
		return nil, err
	}

	lib.d = discovery.New(lib.p, lib.r)

	if err := web.New(e.Group("/lib"), lib.r, lib.p); err != nil {
		return nil, err
	}

	if err := api.New(*e.Group(""), lib.r); err != nil {
		return nil, err
	}

	return lib, nil
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

func (l *Library) ScanAsync() {
	go func() {
		if err := l.Scan(); err != nil {
			slog.Error("Error scanning library", "error", err)
		}
	}()
}
