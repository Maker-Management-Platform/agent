package web

import (
	"log/slog"

	"github.com/eduardooliveira/stLib/v2/library/process"
	"github.com/eduardooliveira/stLib/v2/library/repo"
	"github.com/labstack/echo/v4"
)

type webHandler struct {
	l *slog.Logger
	r *repo.AssetRepo
	p *process.Processor
}

func New(e *echo.Group, r *repo.AssetRepo, p *process.Processor) error {
	wh := &webHandler{
		l: slog.With("module", "library-web"),
		r: r,
		p: p,
	}
	e.GET("", wh.indexHandler)
	e.GET("/details", wh.getAssetDetails)
	e.GET("/edit", wh.editAsset)
	e.POST("/edit", wh.editAsset)
	e.GET("/new", wh.newAsset)
	e.POST("/new", wh.newAsset)
	e.GET("/:assetID", wh.indexHandler)
	e.GET("/list", wh.listHandler)
	e.GET("/:assetID/file", wh.getFileHandler)
	e.GET("/:assetID/extract", wh.extractHandler)
	return nil
}
