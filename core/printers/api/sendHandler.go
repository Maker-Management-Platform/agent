package api

import (
	"log"
	"net/http"

	"github.com/eduardooliveira/stLib/core/data/database"
	"github.com/eduardooliveira/stLib/core/integrations/klipper"
	"github.com/eduardooliveira/stLib/core/integrations/octorpint"
	"github.com/eduardooliveira/stLib/core/state"
	"github.com/labstack/echo/v4"
)

func sendHandler(c echo.Context) error {
	uuid := c.Param("uuid")
	printer, ok := state.Printers[uuid]

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	id := c.Param("id")
	asset, err := database.GetAsset(id)

	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if printer.Type == "klipper" {
		err = klipper.UploadFile(printer, asset)
	}

	if printer.Type == "octoPrint" {
		err = octorpint.UploadFile(printer, asset)
	}

	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
