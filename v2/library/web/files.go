package web

import (
	"errors"
	"net/http"
	"path/filepath"

	"github.com/eduardooliveira/stLib/v2/config"
	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/utils"
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

	if err := h.r.LoadParents(&asset, 1); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if chi.URLParam(r, "download") != "" {
		w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(*asset.Path))
	} else {
		w.Header().Set("Content-Disposition", "inline; filename="+filepath.Base(*asset.Path))
	}

	if utils.VoZ(asset.NodeKind) == entities.NodeKindBundled {
		http.ServeFile(w, r, filepath.Join(config.Cfg.Core.DataFolder, "temp", utils.VoZ(asset.ParentID), filepath.Base(*asset.Path)))
		return
	}
	http.ServeFile(w, r, filepath.Join(*asset.Root, *asset.Path))
}
