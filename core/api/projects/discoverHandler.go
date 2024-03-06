package projects

import (
	"errors"
	"log"
	"net/http"

	"github.com/eduardooliveira/stLib/core/data/database"
	"github.com/eduardooliveira/stLib/core/processing"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func discoverHandler(c echo.Context) error {

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

	_, err = processing.HandlePath(project.FullPath())
	if err != nil {
		log.Printf("error discovering the project %q: %v\n", project.FullPath(), err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
