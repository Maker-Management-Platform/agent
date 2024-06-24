package web

import (
	"github.com/eduardooliveira/stLib/v2/library/web/files"
	"github.com/eduardooliveira/stLib/v2/library/web/roots"
	"github.com/labstack/echo/v4"
)

func Init(e echo.Group) error {
	e.GET("", roots.IndexHandler)
	e.GET("/:assetID", roots.IndexHandler)
	e.GET("/:assetID/file", files.GetFileHandler)
	return nil
}
