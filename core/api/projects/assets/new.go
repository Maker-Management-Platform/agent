package assets

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/eduardooliveira/stLib/core/data/database"
	"github.com/eduardooliveira/stLib/core/downloader/tools"
	"github.com/eduardooliveira/stLib/core/entities"
	"github.com/eduardooliveira/stLib/core/utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func New(c echo.Context) error {

	pAsset := &entities.ProjectAsset{}

	if err := c.Bind(pAsset); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	files := form.File["files"]

	if len(files) == 0 {
		log.Println("No files")
		return c.NoContent(http.StatusBadRequest)
	}

	project, err := database.GetProject(pAsset.ProjectUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	path := utils.ToLibPath(fmt.Sprintf("%s/%s", project.FullPath(), pAsset.Name))

	src, err := files[0].Open()
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	defer src.Close()
	if err = tools.SaveFile(filepath.Join(path, files[0].Filename), src); err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	/*processing.EnqueueInitJob(&processing.ProcessableAsset{
		Name:    files[0].Filename,
		Project: project,
		Origin:  "fs",
	})*/

	return c.NoContent(http.StatusCreated)
}
