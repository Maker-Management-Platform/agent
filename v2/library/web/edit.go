package web

import (
	"net/http"

	"github.com/eduardooliveira/stLib/v2/web"
)

//TODO migrate to chi

func (h webHandler) editAsset(r *http.Request) web.ResponseModel {
	/*
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
		})*/
	return web.ResponseModel{}
}
