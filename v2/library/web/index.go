package web

import (
	"errors"
	"net/http"
	"path"

	"github.com/a-h/templ"
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/eduardooliveira/stLib/v2/config"
	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/library/web/comp"
	"github.com/eduardooliveira/stLib/v2/utils"
	"github.com/eduardooliveira/stLib/v2/web"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func (h webHandler) indexHandler(r *http.Request) web.ResponseModel {
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
		Component: maybeWrapComponent(r, comp.Index(comp.IndexModel{
			Asset: &asset,
			Main:  listComp,
			KindFilter: comp.KindFilterModel{
				Selected:   chi.URLParam(r, "kind"),
				AssetTypes: maputil.Values(config.Cfg.Library.AssetTypes),
			},
			Pagination: *pgModel,
		})),
	}
}

func maybeWrapComponent(r *http.Request, c templ.Component) templ.Component {
	if r.Header.Get("Hx-Boosted") == "true" || r.Header.Get("Hx-Request") != "true" {
		return comp.Wrapper(c)
	}
	return c
}
