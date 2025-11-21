package api

import (
	"log/slog"
	"net/http"

	"github.com/eduardooliveira/stLib/v2/events"
	"github.com/eduardooliveira/stLib/v2/library/process"
	"github.com/eduardooliveira/stLib/v2/library/repo"
	"github.com/go-chi/chi/v5"
)

type APIHandler struct {
	log       *slog.Logger
	repo      repo.AssetRepo
	processor *process.Processor
	eventMgr  *events.EventManager
}

func New(repo repo.AssetRepo, processor *process.Processor) (http.Handler, error) {
	return NewWithEventManager(repo, processor, nil)
}

func NewWithEventManager(repo repo.AssetRepo, processor *process.Processor, eventMgr *events.EventManager) (http.Handler, error) {
	ah := &APIHandler{
		log:       slog.With("module", "library-api"),
		repo:      repo,
		processor: processor,
		eventMgr:  eventMgr,
	}

	r := chi.NewRouter()

	// Asset management endpoints
	r.Get("/", ah.indexHandler)
	r.Get("/{assetID}", ah.indexHandler)
	r.Patch("/{assetID}", ah.patchHandler)
	r.Get("/{assetID}/file", ah.getFileHandler)

	// Downloader endpoints
	r.Post("/download", ah.downloadHandler)
	r.Get("/download/status", ah.downloadStatusHandler)

	return r, nil
}
