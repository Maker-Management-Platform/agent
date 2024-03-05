package projects

import (
	"errors"
	"log"
	"net/http"

	"github.com/eduardooliveira/stLib/core/data/database"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func show(c echo.Context) error {
	uuid := c.Param("uuid")
	rtn, err := database.GetProject(uuid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, rtn)
}
