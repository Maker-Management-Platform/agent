package web

import (
	"errors"
	"net/http"

	"github.com/eduardooliveira/stLib/v2/library/web/comp"
	"github.com/eduardooliveira/stLib/v2/web"
	"gorm.io/gorm"
)

func (h webHandler) viewerListAsset(r *http.Request) web.ResponseModel {
	id := r.URL.Query().Get("assetID")
	if id == "" {
		return web.ResponseModel{
			Status: http.StatusBadRequest,
			Error:  errors.New("Asset ID is required"),
		}
	}
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
		Status:     http.StatusOK,
		Component:  comp.ViewerAssetElement(asset),
		IsFragment: true,
	}
}
