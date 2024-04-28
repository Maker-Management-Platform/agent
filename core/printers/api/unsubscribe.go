package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/eduardooliveira/stLib/core/events"
	"github.com/eduardooliveira/stLib/core/state"
	"github.com/labstack/echo/v4"
)

func unsubscribe(c echo.Context) error {

	session := c.Param("session")
	if session == "" {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("no session provided").Error())
	}

	uuid := c.Param("uuid")
	printer, ok := state.Printers[uuid]

	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, errors.New("printer not found").Error())
	}
	events.UnSubscribe(session, fmt.Sprintf("printer.%s", printer.UUID))

	return c.NoContent(http.StatusOK)
}
