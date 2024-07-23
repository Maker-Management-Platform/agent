package web

import (
	"errors"
	"net/http"

	"github.com/duke-git/lancet/v2/maputil"
	"github.com/eduardooliveira/stLib/v2/config"
	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/library/web/comp"
	"github.com/eduardooliveira/stLib/v2/web"
	corecomp "github.com/eduardooliveira/stLib/v2/web/comp"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (h webHandler) indexHandler(c echo.Context) error {
	roots := config.Cfg.Library.Paths
	var err error
	var asset entities.Asset
	if c.Param("assetID") != "" {
		asset, err = h.r.GetAsset(c.Param("assetID"), true)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return web.Error(c, http.StatusNotFound, "Asset not found")
			}
			return web.Error(c, http.StatusInternalServerError, err.Error())
		}
	} else {
		if len(roots) == 0 {
			return web.Error(c, http.StatusNotFound, "No library paths configured, check library.paths in config.toml")
		}
		asset, err = h.r.GetAssetByRootAndPath(roots[0], ".", true)
	}

	if err != nil {
		return web.Error(c, http.StatusInternalServerError, err.Error())
	}

	err = h.r.LoadParents(&asset, 5, "ID", "Label")
	if err != nil {
		return web.Error(c, http.StatusInternalServerError, err.Error())
	}

	listComp, err := h.list(&listInput{
		c:     c,
		Asset: &asset,
	})
	if err != nil {
		return web.Error(c, http.StatusInternalServerError, err.Error())
	}

	return web.Render(web.ResponseModel{
		Ctx: c,
		S:   http.StatusOK,
		WrapperModel: corecomp.WrapperModel{
			Main: comp.Index(comp.IndexModel{
				Asset: &asset,
				Main:  listComp,
				KindFilter: comp.KindFilterModel{
					Selected:   c.QueryParam("kind"),
					AssetTypes: maputil.Values(config.Cfg.Library.AssetTypes),
				},
			}),
			AsideR: comp.SideBar(),
		},
	})

}
