package roots

import (
	"errors"
	"net/http"

	"github.com/eduardooliveira/stLib/v2/config"
	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/library/svc"
	"github.com/eduardooliveira/stLib/v2/web"
	"github.com/eduardooliveira/stLib/v2/web/components"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func IndexHandler(c echo.Context) error {
	roots := config.Cfg.Library.Paths
	var err error
	var asset *entities.Asset
	if c.Param("assetID") != "" {
		asset, err = svc.GetAsset(c.Param("assetID"), true)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return web.Error(c, http.StatusNotFound, "Asset not found")
			}
			return web.Error(c, http.StatusInternalServerError, err.Error())
		}
	} else {
		asset, err = svc.GetAssetByRootAndPath(roots[0], ".", true)
	}

	if err != nil {
		return web.Error(c, http.StatusInternalServerError, err.Error())
	}

	return web.Render(c, http.StatusOK, components.Wrapper{
		MainContent: Index(asset),
	}, true)
}
