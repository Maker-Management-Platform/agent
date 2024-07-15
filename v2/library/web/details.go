package web

import (
	"errors"
	"net/http"

	"github.com/eduardooliveira/stLib/v2/library/web/comp"
	"github.com/eduardooliveira/stLib/v2/web"
	corecomp "github.com/eduardooliveira/stLib/v2/web/comp"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (h webHandler) getAssetDetails(c echo.Context) error {
	id := c.QueryParam("assetID")
	asset, err := h.r.GetAsset(id, false)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.NoContent(http.StatusNotFound)
		}
		return web.Error(c, http.StatusInternalServerError, err.Error())
	}
	return web.Render(web.ResponseModel{
		Ctx: c,
		S:   http.StatusOK,
		WrapperModel: corecomp.WrapperModel{
			Main: comp.Details(&comp.DetailsModel{
				Asset: &asset,
			}),
		},
		IsFragment: true,
	})
}
