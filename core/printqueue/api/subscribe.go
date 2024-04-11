package printqueue

import (
	"errors"
	"net/http"

	"github.com/eduardooliveira/stLib/core/events"
	"github.com/eduardooliveira/stLib/core/printqueue/state"
	"github.com/labstack/echo/v4"
)

func subscribe(c echo.Context) error {

	session := c.Param("session")
	if session == "" {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("no session provided").Error())
	}

	err := events.Subscribe(session, "printQueue", state.GetEventPublisher())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
