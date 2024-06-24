package files

import (
	"errors"
	"net/http"
	"path/filepath"

	"github.com/eduardooliveira/stLib/v2/library/svc"
	"github.com/eduardooliveira/stLib/v2/web"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetFileHandler(c echo.Context) error {
	id := c.Param("assetID")
	asset, err := svc.GetAsset(id, false)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.NoContent(http.StatusNotFound)
		}
		return web.Error(c, http.StatusInternalServerError, err.Error())
	}
	if c.QueryParam("download") != "" {
		return c.Attachment(filepath.Join(*asset.Root, *asset.Path), "")
	}
	return c.Inline(filepath.Join(*asset.Root, *asset.Path), "")

}
