package web

import (
	"errors"
	"net/http"

	"github.com/eduardooliveira/stLib/v2/config"
	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/library/web/comp"
	"github.com/eduardooliveira/stLib/v2/web"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (h webHandler) indexHandler(c echo.Context) error {
	roots := config.Cfg.Library.Paths
	var err error
	var asset *entities.Asset
	if c.Param("assetID") != "" {
		asset, err = h.r.GetAsset(c.Param("assetID"), true)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return web.Error(c, http.StatusNotFound, "Asset not found")
			}
			return web.Error(c, http.StatusInternalServerError, err.Error())
		}
	} else {
		asset, err = h.r.GetAssetByRootAndPath(roots[0], ".", true)
	}

	if err != nil {
		return web.Error(c, http.StatusInternalServerError, err.Error())
	}

	err = h.r.LoadParents(asset, 5, "ID", "Label")
	if err != nil {
		return web.Error(c, http.StatusInternalServerError, err.Error())
	}

	listComp, err := h.list(&listInput{
		c:     c,
		Asset: asset,
	})
	if err != nil {
		return web.Error(c, http.StatusInternalServerError, err.Error())
	}

	return web.Render(web.ResponseModel{
		Ctx: c,
		S:   http.StatusOK,
		MainComponent: comp.Index(comp.IndexModel{
			Asset: asset,
			Main:  listComp,
		}),
	})

}
