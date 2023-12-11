package tempfiles

import (
	"net/http"

	"github.com/eduardooliveira/stLib/core/models"
	"github.com/eduardooliveira/stLib/core/state"
	"github.com/labstack/echo/v4"
)

func index(c echo.Context) error {
	rtn := make([]*models.TempFile, 0)
	for _, t := range state.TempFiles {
		rtn = append(rtn, t)
	}
	return c.JSON(http.StatusOK, rtn)
}

func move(c echo.Context) error {
	rtn := make([]*models.TempFile, 0)
	for _, t := range state.TempFiles {
		rtn = append(rtn, t)
	}
	return c.JSON(http.StatusOK, rtn)
}
