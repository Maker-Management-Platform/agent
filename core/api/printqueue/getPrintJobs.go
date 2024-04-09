package printqueue

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/eduardooliveira/stLib/core/data/database"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func getPrintJobs(c echo.Context) error {

	if c.QueryParam("states") == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "states parameter is required")

	}
	states := strings.Split(c.QueryParam("states"), ",")

	rtn, err := database.GetPrintJobs(states)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, rtn)
}
