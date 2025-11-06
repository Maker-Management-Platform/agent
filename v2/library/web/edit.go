package web

import (
	"errors"
	"net/http"

	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/library/web/comp"
	"github.com/eduardooliveira/stLib/v2/web"
	"github.com/ggicci/httpin"
	"gorm.io/gorm"
)

//TODO migrate to chi

func (h webHandler) getEditAssetHandler(r *http.Request) web.ResponseModel {
	id := r.URL.Query().Get("assetID")
	asset, err := h.r.GetAsset(id, false)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return web.ResponseModel{
				Status: http.StatusNotFound,
				Error:  err,
			}
		}
		return web.ResponseModel{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}
	return web.ResponseModel{
		Status: http.StatusOK,
		Component: comp.Edit(&comp.EditModel{
			Asset: &asset,
		}),
		IsFragment: true,
	}
}

func (h webHandler) postEditAssetHandler(r *http.Request) web.ResponseModel {

	asset := r.Context().Value(httpin.Input).(*entities.Asset)
	err := h.r.UpdateAsset(asset)
	if err != nil {
		return web.ResponseModel{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}
	return web.ResponseModel{
		Status: http.StatusOK,
		Component: comp.Edit(&comp.EditModel{
			Asset: asset,
		}),
		IsFragment: true,
	}
}

/*
func (h webHandler) editAsset(r *http.Request) web.ResponseModel {

	var asset entities.Asset
	var err error
	var action string

	if c.Request().Method == http.MethodGet {
		id := c.QueryParam("assetID")
		asset, err = h.r.GetAsset(id, false)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.NoContent(http.StatusNotFound)
			}
			return web.Error(c, http.StatusInternalServerError, err.Error())
		}
		action = "edit"
	} else if c.Request().Method == http.MethodPost {
		err = c.Bind(&asset)
		if err != nil {
			return web.Error(c, http.StatusBadRequest, err.Error())
		}
		err = h.r.UpdateAsset(&asset)
		if err != nil {
			return web.Error(c, http.StatusInternalServerError, err.Error())
		}
		action = "save"
	}
	return web.Render(web.ResponseModel{
		Ctx: c,
		S:   http.StatusOK,
		WrapperModel: corecomp.WrapperModel{
			Main: comp.Edit(&comp.EditModel{
				Asset:  &asset,
				Action: action,
			}),
		},
		IsFragment: true,
	})
	return web.ResponseModel{}
}
*/
