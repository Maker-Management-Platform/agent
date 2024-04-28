package api

import (
	"log"
	"net/http"

	"github.com/eduardooliveira/stLib/core/entities"
	"github.com/eduardooliveira/stLib/core/state"
	"github.com/labstack/echo/v4"
)

func edit(c echo.Context) error {
	pPrinter := &entities.Printer{}

	if err := c.Bind(pPrinter); err != nil {
		log.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}

	if pPrinter.UUID != c.Param("uuid") {
		return c.NoContent(http.StatusBadRequest)
	}

	printer, ok := state.Printers[pPrinter.UUID]

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	printer.Name = pPrinter.Name
	printer.Address = pPrinter.Address
	printer.Type = pPrinter.Type
	printer.ApiKey = pPrinter.ApiKey

	state.Printers[printer.UUID] = printer
	state.PersistPrinters()

	return c.JSON(http.StatusCreated, state.Printers[printer.UUID])
}
