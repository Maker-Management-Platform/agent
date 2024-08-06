package web

import (
	"errors"
	"net/http"
	"path"

	"github.com/duke-git/lancet/v2/maputil"
	"github.com/eduardooliveira/stLib/v2/config"
	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/library/web/comp"
	"github.com/eduardooliveira/stLib/v2/utils"
	"github.com/eduardooliveira/stLib/v2/web"
	corecomp "github.com/eduardooliveira/stLib/v2/web/comp"
	"github.com/go-chi/chi/v5"
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

	listComp, _, err := h.list(&listInput{
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

func (h webHandler) indexHandlerChi(r *http.Request) web.ResponseModel {
	var err error
	var asset entities.Asset
	if chi.URLParam(r, "assetID") != "" {
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
	} else {
		if len(config.Cfg.Library.FileSystems) == 0 || (len(config.Cfg.Library.FileSystems) == 1 && config.Cfg.Library.FileSystems[0].Path == "change_me") {
			return web.ResponseModel{
				S:     http.StatusNotFound,
				Error: errors.New("No library paths configured, check library.paths in config.toml"),
			}
		}
		roots, err := h.r.GetAssetRoots(true)
		if err != nil {
			return web.ResponseModel{
				S:     http.StatusInternalServerError,
				Error: err,
			}
		}
		if len(roots) > 1 {
			asset = entities.Asset{
				ID:           "",
				Label:        utils.Ptr("Libraries"),
				NestedAssets: roots,
			}
		} else {
			asset = utils.VoZ(roots[0])
		}
	}

	if err != nil {
		return web.ResponseModel{
			S:     http.StatusInternalServerError,
			Error: err,
		}
	}

	if asset.ID != "" {

		err = h.r.LoadParents(&asset, 5, "ID", "Label")
		if err != nil {
			return web.ResponseModel{
				S:     http.StatusInternalServerError,
				Error: err,
			}
		}
	}

	listComp, pgModel, err := h.list(&listInput{
		Asset: &asset,
		r:     r,
	})
	if err != nil {
		return web.ResponseModel{
			S:     http.StatusInternalServerError,
			Error: err,
		}
	}
	return web.ResponseModel{
		S:         http.StatusOK,
		PushState: path.Join("/lib", asset.ID),
		Component: comp.Index(comp.IndexModel{
			Asset: &asset,
			Main:  listComp,
			KindFilter: comp.KindFilterModel{
				Selected:   chi.URLParam(r, "kind"),
				AssetTypes: maputil.Values(config.Cfg.Library.AssetTypes),
			},
			Pagination: *pgModel,
		}),
	}

}
