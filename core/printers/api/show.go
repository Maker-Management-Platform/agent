package api

import (
	"net/http"

	"github.com/eduardooliveira/stLib/core/state"
	"github.com/labstack/echo/v4"
)

func show(c echo.Context) error {
	uuid := c.Param("uuid")
	printer, ok := state.Printers[uuid]
	if !ok {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, printer)
}
