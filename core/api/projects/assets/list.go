package assets

import (
	"errors"
	"log"
	"net/http"

	"github.com/eduardooliveira/stLib/core/data/database"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func List(c echo.Context) error {
	uuid := c.Param("uuid")
	if uuid == "" {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("missing project uuid"))
	}
	rtn, err := database.GetAssetsByProject(uuid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, rtn)
}
