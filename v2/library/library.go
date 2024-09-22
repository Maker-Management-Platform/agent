package library

import (
	"log/slog"
	"net/http"

	"github.com/eduardooliveira/stLib/v2/config"
	"github.com/eduardooliveira/stLib/v2/library/discovery"
	"github.com/eduardooliveira/stLib/v2/library/libfs"
	"github.com/eduardooliveira/stLib/v2/library/process"
	"github.com/eduardooliveira/stLib/v2/library/repo"
	"github.com/eduardooliveira/stLib/v2/library/web"
	"golang.org/x/net/context"
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

	if err = libfs.LoadFSs(); err != nil {
		return nil, nil, nil, err
	}

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

func (l Library) ScanFS(ctx context.Context) error {
	eg := errgroup.Group{}
	if len(config.Cfg.Library.FileSystems) == 0 || (len(config.Cfg.Library.FileSystems) == 1 && config.Cfg.Library.FileSystems[0].Path == "change_me") {
		slog.Warn("invalid library file systems configured")
		return nil
	}
	for _, ffs := range libfs.GetFSs() {
		eg.Go(l.d.DiscoverFS(ctx, ffs).Run)
	}
	return eg.Wait()
}

func (l *Library) ScanAsync(ctx context.Context) {
	go func() {
		if err := l.ScanFS(ctx); err != nil {
			slog.Error("Error scanning library", "error", err)
		}
	}()
}
