package printqueue

import (
	"log"
	"net/http"

	"github.com/eduardooliveira/stLib/core/printqueue/state"
	"github.com/labstack/echo/v4"
)

func validate(c echo.Context) error {

	if c.QueryParam("jobUuid") == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "job uuid parameter is required")
	}

	if c.QueryParam("state") == "ok" || c.QueryParam("state") == "nok" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid state parameter")
	}

	if _, err := state.Validate(c.QueryParam("jobUuid"), c.QueryParam("state")); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
