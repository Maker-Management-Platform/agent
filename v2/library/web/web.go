package web

import (
	"log/slog"
	"net/http"

	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/library/process"
	"github.com/eduardooliveira/stLib/v2/library/repo"
	"github.com/eduardooliveira/stLib/v2/library/web/comp"
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
	r.With(httpin.NewInput(comp.SearchModel{})).Get("/", web.R(wh.indexHandler))
	r.Get("/{assetID}", web.R(wh.indexHandler))

	//r.Get("/sidebar", web.R(web.RenderFragment(comp.Sidebar())))
	r.Get("/viewer/list", web.R(wh.viewerListAsset))
	r.Get("/list", web.R(wh.listHandler))
	r.Get("/{assetID}/file", wh.getFileHandler)
	r.Get("/details", web.R(wh.getAssetDetailsHandler))

	r.Get("/edit", web.R(wh.getEditAssetHandler))
	r.With(httpin.NewInput(entities.Asset{})).Post("/edit", web.R(wh.postEditAssetHandler))

	r.Get("/new", web.R(wh.newAssetHandler))
	r.With(httpin.NewInput(newAssetRequest{})).Post("/new", web.R(wh.newAssetHandler))

	return r, nil
}
