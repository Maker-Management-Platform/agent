package printqueue

import (
	"log"
	"net/http"
	"strconv"

	"github.com/eduardooliveira/stLib/core/printqueue/state"
	"github.com/labstack/echo/v4"
)

func move(c echo.Context) error {

	if c.QueryParam("jobUuid") == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "job uuid parameter is required")
	}

	if c.QueryParam("position") == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "position parameter is required")
	}

	pos, err := strconv.Atoi(c.QueryParam("position"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "position parameter must be an integer")
	}

	if _, err := state.MovePrintJob(c.QueryParam("jobUuid"), pos); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
