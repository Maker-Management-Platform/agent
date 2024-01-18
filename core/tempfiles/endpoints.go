package tempfiles

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/eduardooliveira/stLib/core/data/database"
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

	project, err := database.GetProject(tempFile.ProjectUUID)

	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	dst := utils.ToLibPath(fmt.Sprintf("%s/%s", project.FullPath(), tempFile.Name))

	err = utils.Move(filepath.Clean(fmt.Sprintf("temp/%s", tempFile.Name)), dst, false)

	if err != nil {
		log.Println("Error moving temp file: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	f, err := os.Open(dst)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	asset, nestedAssets, err := models.NewProjectAsset(tempFile.Name, project, f)

	if err != nil {
		log.Println("error initializing asset: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if err = database.InsertAsset(asset); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	for _, a := range nestedAssets {
		if project.DefaultImageID == "" && a.AssetType == "image" {
			project.DefaultImageID = a.ID
			if err := database.UpdateProject(project); err != nil {
				log.Println(err)
			}
		}

		err := database.InsertAsset(a)
		if err != nil {
			log.Println(err)
		}
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
