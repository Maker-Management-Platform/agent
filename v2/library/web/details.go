package web

import (
	"errors"
	"net/http"

	"github.com/eduardooliveira/stLib/v2/library/web/comp"
	"github.com/eduardooliveira/stLib/v2/web"
	"gorm.io/gorm"
)

func (h webHandler) getAssetDetailsHandler(r *http.Request) web.ResponseModel {
	id := r.URL.Query().Get("assetID")
	asset, err := h.r.GetAsset(id, false)
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
	return web.ResponseModel{
		S: http.StatusOK,
		Component: comp.Details(&comp.DetailsModel{
			Asset: &asset,
		}),
		IsFragment: true,
	}
}
