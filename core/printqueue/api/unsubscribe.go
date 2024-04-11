package printqueue

import (
	"errors"
	"net/http"

	"github.com/eduardooliveira/stLib/core/events"
	"github.com/labstack/echo/v4"
)

func unSubscribe(c echo.Context) error {

	session := c.Param("session")
	if session == "" {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("no session provided").Error())
	}

	events.UnSubscribe(session, "printQueue")

	return c.NoContent(http.StatusOK)
}
