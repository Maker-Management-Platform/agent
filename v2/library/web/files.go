package web

import (
	"errors"
	"net/http"
	"path/filepath"

	"github.com/eduardooliveira/stLib/v2/web"
	"github.com/go-chi/chi/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (h webHandler) getFileHandler(c echo.Context) error {
	id := c.Param("assetID")
	asset, err := h.r.GetAsset(id, false)
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
func (h webHandler) getFileHandlerChi(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "assetID")
	if id == "" {
		http.Error(w, "Asset ID is required", http.StatusBadRequest)
		return
	}
	asset, err := h.r.GetAsset(id, false)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if chi.URLParam(r, "download") != "" {
		w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(*asset.Path))
	} else {
		w.Header().Set("Content-Disposition", "inline; filename="+filepath.Base(*asset.Path))
	}

	http.ServeFile(w, r, filepath.Join(*asset.Root, *asset.Path))
}
