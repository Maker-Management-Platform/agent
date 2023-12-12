package printers

import (
	"net/http"

	"github.com/eduardooliveira/stLib/core/models"
	"github.com/labstack/echo/v4"
)

func save(c echo.Context) error {

	pPrinter := &models.Printer{}
	err := c.Bind(pPrinter)

	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if pPrinter.UUID != "" {

	}

	return c.NoContent(http.StatusInternalServerError)
}
