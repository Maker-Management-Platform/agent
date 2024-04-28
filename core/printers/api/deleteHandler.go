package api

import (
	"net/http"

	"github.com/eduardooliveira/stLib/core/state"
	"github.com/labstack/echo/v4"
)

func deleteHandler(c echo.Context) error {
	uuid := c.Param("uuid")
	printer, ok := state.Printers[uuid]

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	delete(state.Printers, printer.UUID)

	err := state.PersistPrinters()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, printer)
}
