package web

import (
	"log/slog"
	"net/http"

	"github.com/eduardooliveira/stLib/v2/library/process"
	"github.com/eduardooliveira/stLib/v2/library/repo"
	"github.com/eduardooliveira/stLib/v2/web"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
)

type webHandler struct {
	l *slog.Logger
	r *repo.AssetRepo
	p *process.Processor
}

func New(repo *repo.AssetRepo, p *process.Processor) (http.Handler, error) {
	wh := &webHandler{
		l: slog.With("module", "library-web"),
		r: repo,
		p: p,
	}
	r := chi.NewRouter()
	r.Get("/", web.R(wh.indexHandler))
	r.Get("/{assetID}", web.R(wh.indexHandler))

	r.Get("/sidebar", web.R(wh.sidebarHandler))
	r.Get("/list", web.R(wh.listHandler))
	r.Get("/{assetID}/file", wh.getFileHandler)
	r.Get("/details", web.R(wh.getAssetDetailsHandler))

	r.Get("/new", web.R(wh.newAssetHandler))
	r.With(httpin.NewInput(newAssetRequest{})).Post("/new", web.R(wh.newAssetHandler))

	return r, nil
}
