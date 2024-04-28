package api

import (
	"net/http"

	"github.com/eduardooliveira/stLib/core/entities"
	"github.com/eduardooliveira/stLib/core/state"
	"github.com/labstack/echo/v4"
)

func new(c echo.Context) error {

	pPrinter := &entities.Printer{}
	err := c.Bind(pPrinter)

	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	printer := entities.NewPrinter()
	printer.Name = pPrinter.Name
	printer.Address = pPrinter.Address
	printer.Type = pPrinter.Type
	printer.ApiKey = pPrinter.ApiKey

	state.Printers[printer.UUID] = printer
	state.PersistPrinters()

	return c.JSON(http.StatusCreated, state.Printers[printer.UUID])
}
