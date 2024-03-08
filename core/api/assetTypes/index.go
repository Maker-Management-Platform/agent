package assettypes

import (
	"net/http"

	"github.com/duke-git/lancet/v2/maputil"
	"github.com/eduardooliveira/stLib/core/entities"
	"github.com/eduardooliveira/stLib/core/state"
	"github.com/labstack/echo/v4"
)

func index(c echo.Context) error {
	return c.JSON(http.StatusOK, maputil.Values[string, *entities.AssetType](state.AssetTypes))
}
