package projects

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/eduardooliveira/stLib/core/data/database"
	"github.com/eduardooliveira/stLib/core/utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func deleteHandler(c echo.Context) error {

	uuid := c.Param("uuid")

	if uuid == "" {
		return c.NoContent(http.StatusBadRequest)
	}
	project, err := database.GetProject(uuid)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = os.RemoveAll(utils.ToLibPath(project.FullPath()))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	err = utils.DeleteAssetsFolder(project.UUID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := database.DeleteProject(uuid); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}
