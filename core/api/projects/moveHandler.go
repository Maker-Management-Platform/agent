package projects

import (
	"errors"
	"log"
	"net/http"
	"path/filepath"

	"github.com/eduardooliveira/stLib/core/data/database"
	"github.com/eduardooliveira/stLib/core/entities"
	"github.com/eduardooliveira/stLib/core/utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func moveHandler(c echo.Context) error {
	pproject := &entities.Project{}

	if err := c.Bind(pproject); err != nil {
		log.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}

	if pproject.UUID != c.Param("uuid") {
		return c.NoContent(http.StatusBadRequest)
	}

	project, err := database.GetProject(pproject.UUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	pproject.Path = filepath.Clean(pproject.Path)
	pproject.Name = project.Name
	err = utils.Move(project.FullPath(), pproject.FullPath(), true)

	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	project.Path = filepath.Clean(pproject.Path)

	err = database.UpdateProject(project)

	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, struct {
		UUID string `json:"uuid"`
		Path string `json:"path"`
	}{project.UUID, project.Path})
}
