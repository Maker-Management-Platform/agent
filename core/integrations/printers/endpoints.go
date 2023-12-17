package printers

import (
	"log"
	"net/http"

	"github.com/eduardooliveira/stLib/core/integrations/klipper"
	"github.com/eduardooliveira/stLib/core/models"
	"github.com/eduardooliveira/stLib/core/state"
	"github.com/labstack/echo/v4"
)

func index(c echo.Context) error {
	rtn := make([]*models.Printer, 0)
	for _, p := range state.Printers {
		rtn = append(rtn, p)
	}
	return c.JSON(http.StatusOK, rtn)
}

func show(c echo.Context) error {
	uuid := c.Param("uuid")
	printer, ok := state.Printers[uuid]

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, printer)
}

func deleteHandler(c echo.Context) error {
	uuid := c.Param("uuid")
	printer, ok := state.Printers[uuid]

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	delete(state.Printers, printer.UUID)

	err := state.PercistPrinters()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, printer)
}

func new(c echo.Context) error {

	pPrinter := &models.Printer{}
	err := c.Bind(pPrinter)

	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	printer := models.NewPrinter()
	printer.Name = pPrinter.Name
	printer.Address = pPrinter.Address
	printer.Type = pPrinter.Type

	state.Printers[printer.UUID] = printer
	state.PercistPrinters()

	return c.JSON(http.StatusCreated, state.Printers[printer.UUID])
}

func edit(c echo.Context) error {
	pPrinter := &models.Printer{}

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

	state.Printers[printer.UUID] = printer
	state.PercistPrinters()

	return c.JSON(http.StatusCreated, state.Printers[printer.UUID])
}

func testConnection(c echo.Context) error {

	pPrinter := &models.Printer{}
	err := c.Bind(pPrinter)

	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if pPrinter.Type == "klipper" {
		err = klipper.ConntectionStatus(pPrinter)
		log.Println(err)
		return c.JSON(http.StatusOK, pPrinter)
	}

	return c.NoContent(http.StatusInternalServerError)
}
