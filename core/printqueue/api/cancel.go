package printqueue

import (
	"errors"
	"log"
	"net/http"

	"github.com/eduardooliveira/stLib/core/printqueue/state"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func cancel(c echo.Context) error {

	uuid := c.Param("uuid")

	if uuid == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	jobs, err := state.ChangePrintJobState(uuid, "canceled")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, jobs)
}
