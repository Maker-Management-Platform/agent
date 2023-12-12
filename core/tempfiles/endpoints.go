package tempfiles

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/eduardooliveira/stLib/core/discovery"
	"github.com/eduardooliveira/stLib/core/models"
	"github.com/eduardooliveira/stLib/core/state"
	"github.com/eduardooliveira/stLib/core/utils"
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
	uuid := c.Param("uuid")

	if uuid == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	_, ok := state.TempFiles[uuid]

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	tempFile := &models.TempFile{}

	if err := c.Bind(tempFile); err != nil {
		log.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}

	if uuid != tempFile.UUID {
		return c.NoContent(http.StatusBadRequest)
	}

	project, ok := state.Projects[tempFile.ProjectUUID]

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	err := os.Rename(filepath.Clean(fmt.Sprintf("temp/%s", tempFile.Name)), fmt.Sprintf("%s/%s", utils.ToLibPath(project.FullPath()), tempFile.Name))

	if err != nil {
		log.Println("Error moving temp file: ", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	err = discovery.DiscoverProjectAssets(project)

	if err != nil {
		log.Println("Error discovering project assets: ", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	delete(state.TempFiles, uuid)
	return c.NoContent(http.StatusOK)
}

func deleteTempFile(c echo.Context) error {

	uuid := c.Param("uuid")

	if uuid == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	tempFile, ok := state.TempFiles[uuid]

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	err := os.Remove(fmt.Sprintf("temp/%s", tempFile.Name))
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	delete(state.TempFiles, uuid)

	return c.NoContent(http.StatusOK)
}
