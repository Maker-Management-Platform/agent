package tags

import (
	"log"
	"net/http"

	"github.com/eduardooliveira/stLib/core/data/database"
	"github.com/labstack/echo/v4"
)

func index(c echo.Context) error {
	rtn, err := database.GetTags()
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, rtn)
}
