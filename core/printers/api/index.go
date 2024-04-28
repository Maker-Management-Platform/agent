package api

import (
	"net/http"

	"github.com/eduardooliveira/stLib/core/entities"
	"github.com/eduardooliveira/stLib/core/state"
	"github.com/labstack/echo/v4"
)

func index(c echo.Context) error {
	rtn := make([]*entities.Printer, 0)
	for _, p := range state.Printers {
		rtn = append(rtn, p)
	}
	return c.JSON(http.StatusOK, rtn)
}
