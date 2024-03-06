package projects

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/eduardooliveira/stLib/core/data/database"
	"github.com/eduardooliveira/stLib/core/entities"
	"github.com/eduardooliveira/stLib/core/utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func save(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}

	projectPayload := form.Value["payload"]
	if len(projectPayload) != 1 {
		log.Println("more payloads than expected")
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("more payloads than expected"))
	}

	pproject := &entities.Project{}

	err = json.Unmarshal([]byte(projectPayload[0]), pproject)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Bind(pproject); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if pproject.UUID != c.Param("uuid") {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("parameter mismatch"))
	}

	project, err := database.GetProject(c.Param("uuid"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if pproject.Name != project.Name {

		err := utils.Move(project.FullPath(), pproject.FullPath(), true)

		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	err = database.UpdateProject(pproject)

	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, pproject)
}
