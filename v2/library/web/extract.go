package web

import (
	"net/http"

	"github.com/eduardooliveira/stLib/v2/library/web/comp"
	"github.com/eduardooliveira/stLib/v2/web"
	corecomp "github.com/eduardooliveira/stLib/v2/web/comp"
	"github.com/labstack/echo/v4"
)

func (h webHandler) extractHandler(c echo.Context) error {
	if c.Param("assetID") == "" {
		return web.Error(c, http.StatusBadRequest, "assetID is required")
	}

	asset, err := h.r.GetAsset(c.Param("assetID"), false)
	if err != nil {
		return web.Error(c, http.StatusInternalServerError, err.Error())
	}

	if err := h.r.LoadParents(&asset, 1); err != nil {
		return web.Error(c, http.StatusInternalServerError, err.Error())
	}
	p, err := h.p.ProcessBundled(&asset)
	if err != nil {
		return web.Error(c, http.StatusInternalServerError, err.Error())
	}
	if err := p.Wait(); err != nil {
		return web.Error(c, http.StatusInternalServerError, err.Error())
	}

	return web.Render(web.ResponseModel{
		Ctx: c,
		S:   http.StatusOK,
		WrapperModel: corecomp.WrapperModel{
			Main: comp.AssetCard(&comp.AssetCardModel{
				Asset: &asset,
			}),
		},
		IsFragment: true,
	})
}
