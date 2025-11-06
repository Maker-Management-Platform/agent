package api

import (
	"log/slog"
	"net/http"

	"github.com/eduardooliveira/stLib/v2/library/process"
	"github.com/eduardooliveira/stLib/v2/library/repo"
	"github.com/go-chi/chi/v5"
)

type APIHandler struct {
	log       *slog.Logger
	repo      repo.AssetRepo
	processor *process.Processor
}

func New(repo repo.AssetRepo, processor *process.Processor) (http.Handler, error) {
	ah := &APIHandler{
		log:       slog.With("module", "library-api"),
		repo:      repo,
		processor: processor,
	}
	_ = ah
	r := chi.NewRouter()
	r.Get("/", ah.indexHandler)
	r.Get("/{assetID}", ah.indexHandler)
	r.Patch("/{assetID}", ah.patchHandler)
	r.Get("/{assetID}/file", ah.getFileHandler)
	return r, nil
}
