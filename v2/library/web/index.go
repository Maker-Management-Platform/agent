package web

import (
	"errors"
	"net/http"
	"net/url"
	"path"

	"github.com/a-h/templ"
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/eduardooliveira/stLib/v2/config"
	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/library/web/comp"
	"github.com/eduardooliveira/stLib/v2/utils"
	"github.com/eduardooliveira/stLib/v2/web"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func (h webHandler) indexHandler(r *http.Request) web.ResponseModel {
	var asset entities.Asset
	var resp *web.ResponseModel

	if r.URL.Query().Has("search") {
		return doSearch(h, r)
	} else if chi.URLParam(r, "assetID") != "" {
		asset, resp = withAssetID(h, r, chi.URLParam(r, "assetID"))
		if resp != nil {
			return *resp
		}
	} else {
		asset, resp = roots(h, r)
		if resp != nil {
			return *resp
		}
	}

	if asset.ID != "" {

		err := h.r.LoadParents(&asset, 5, "ID", "Label")
		if err != nil {
			return web.ResponseModel{
				Status: http.StatusInternalServerError,
				Error:  err,
			}
		}
	}

	listComp, pgModel, err := h.list(&listInput{
		Asset: &asset,
		r:     r,
	})
	if err != nil {
		return web.ResponseModel{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}

	return web.ResponseModel{
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

func withAssetID(h webHandler, r *http.Request, assetID string) (entities.Asset, *web.ResponseModel) {
	asset, err := h.r.GetAsset(assetID, true)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Asset{}, &web.ResponseModel{
				Status: http.StatusNotFound,
				Error:  errors.New("Asset not found"),
			}
		}
		return entities.Asset{}, &web.ResponseModel{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}
	return asset, nil
}

func doSearch(h webHandler, r *http.Request) web.ResponseModel {
	pushParams := r.URL.Query()
	pushURL := url.URL{
		Path:     "/lib",
		RawQuery: pushParams.Encode(),
	}
	req := r.Context().Value(httpin.Input).(*comp.SearchModel)

	if req == nil || !req.Search {
		return web.ResponseModel{
			Status: http.StatusBadRequest,
			Error:  errors.New("missing search query"),
		}
	}
	assets, err := h.r.SearchAsset(req.Name, req.Tags)
	if err != nil {
		return web.ResponseModel{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}
	asset := &entities.Asset{
		ID:           "",
		Label:        utils.Ptr("Search results"),
		NestedAssets: assets,
	}
	return web.ResponseModel{
		PushState: pushURL.String(),
		Component: maybeWrapComponent(r, comp.Index(comp.IndexModel{
			Asset:       asset,
			SearchModel: *req,
			Main:        comp.List(comp.ListModel{Asset: asset}),
		})),
	}
}

func roots(h webHandler, r *http.Request) (entities.Asset, *web.ResponseModel) {
	if len(config.Cfg.Library.FileSystems) == 0 || (len(config.Cfg.Library.FileSystems) == 1 && config.Cfg.Library.FileSystems[0].Path == "change_me") {
		return entities.Asset{}, &web.ResponseModel{
			Status: http.StatusNotFound,
			Error:  errors.New("No library paths configured, check library.paths in config.toml"),
		}
	}
	roots, err := h.r.GetAssetRoots(true)
	if err != nil {
		return entities.Asset{}, &web.ResponseModel{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}
	var asset entities.Asset
	if len(roots) > 1 {
		asset = entities.Asset{
			ID:           "",
			Label:        utils.Ptr("Libraries"),
			NestedAssets: roots,
		}
	} else {
		asset = utils.VoZ(roots[0])
	}
	return asset, nil
}
