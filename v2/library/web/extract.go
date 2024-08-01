package web

import (
	"errors"
	"net/http"

	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/library/web/comp"
	"github.com/eduardooliveira/stLib/v2/web"
	corecomp "github.com/eduardooliveira/stLib/v2/web/comp"
	"github.com/go-chi/chi/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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

func (h webHandler) extractHandlerChi(r *http.Request) web.ResponseModel {
	var err error
	var asset entities.Asset
	if chi.URLParam(r, "assetID") == "" {
		return web.ResponseModel{
			S:     http.StatusBadRequest,
			Error: errors.New("Asset ID is required"),
		}
	}

	asset, err = h.r.GetAsset(chi.URLParam(r, "assetID"), true)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return web.ResponseModel{
				S:     http.StatusNotFound,
				Error: err,
			}
		}
		return web.ResponseModel{
			S:     http.StatusInternalServerError,
			Error: err,
		}
	}

	if err := h.r.LoadParents(&asset, 1); err != nil {
		return web.ResponseModel{
			S:     http.StatusInternalServerError,
			Error: err,
		}
	}

	p, err := h.p.ProcessBundled(&asset)
	if err != nil {
		return web.ResponseModel{
			S:     http.StatusInternalServerError,
			Error: err,
		}
	}
	if err := p.Wait(); err != nil {
		return web.ResponseModel{
			S:     http.StatusInternalServerError,
			Error: err,
		}
	}
	return web.ResponseModel{
		S: http.StatusOK,
		Component: comp.AssetCard(&comp.AssetCardModel{
			Asset: &asset,
		}),
		IsFragment: true,
	}

}
