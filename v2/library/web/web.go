package web

import (
	"log/slog"

	"github.com/eduardooliveira/stLib/v2/library/repo"
	"github.com/labstack/echo/v4"
)

type webHandler struct {
	l *slog.Logger
	r *repo.AssetRepo
}

func New(e *echo.Group, r *repo.AssetRepo) error {
	wh := &webHandler{
		l: slog.With("module", "library-web"),
		r: r,
	}
	e.GET("", wh.indexHandler)
	e.GET("/:assetID", wh.indexHandler)
	e.GET("/:assetID/list", wh.listHandler)
	e.GET("/:assetID/file", wh.getFileHandler)
	return nil
}
