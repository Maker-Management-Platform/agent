package api

import (
	"log"
	"net/http"

	"github.com/eduardooliveira/stLib/core/entities"
	"github.com/eduardooliveira/stLib/core/integrations/klipper"
	"github.com/eduardooliveira/stLib/core/integrations/octorpint"
	"github.com/labstack/echo/v4"
)

func testConnection(c echo.Context) error {

	pPrinter := &entities.Printer{}
	err := c.Bind(pPrinter)

	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if pPrinter.Type == "klipper" {
		err = klipper.ConnectionStatus(pPrinter)
		log.Println(err)
		return c.JSON(http.StatusOK, pPrinter)
	} else if pPrinter.Type == "octoPrint" {
		err = octorpint.ConnectionStatus(pPrinter)
		log.Println(err)
		return c.JSON(http.StatusOK, pPrinter)
	}

	return c.NoContent(http.StatusInternalServerError)
}
