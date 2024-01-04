package assets

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/eduardooliveira/stLib/core/data/database"
	"github.com/eduardooliveira/stLib/core/models"
	"github.com/eduardooliveira/stLib/core/utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

/* not in use, needs to be migrated to database
func save(c echo.Context) error {
	sha1 := c.Param("sha1")

	if sha1 == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	asset, ok := state.Assets[sha1]

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}
	project, ok := state.Projects[asset.ProjectUUID]

	if !ok {
		return c.NoContent(http.StatusInternalServerError)
	}

	pAsset := &models.ProjectAsset{}
	err := c.Bind(pAsset)

	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	oldPath := utils.ToLibPath(fmt.Sprintf("%s/%s", project.FullPath(), asset.Name))

	if pAsset.ProjectUUID != asset.ProjectUUID {

		newProject, ok := state.Projects[pAsset.ProjectUUID]

		if !ok {
			return c.NoContent(http.StatusNotFound)
		}

		newPath := utils.ToLibPath(fmt.Sprintf("%s/%s", newProject.FullPath(), pAsset.Name))
		err = utils.Move(oldPath, newPath)

		if err != nil {
			log.Println("move", err)
			return c.NoContent(http.StatusInternalServerError)
		}

		delete(state.Assets, sha1)
		delete(project.Assets, sha1)

		f, err := os.Open(newPath)
		if err != nil {
			log.Println("open", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		defer f.Close()

		asset, _, err := models.NewProjectAsset(pAsset.Name, newProject, f)

		if err != nil {
			log.Println("new", err)
			return c.NoContent(http.StatusInternalServerError)
		}

		newProject.Assets[asset.SHA1] = asset
		state.Assets[asset.SHA1] = asset
	}

	if pAsset.Name != asset.Name {
		newPath := utils.ToLibPath(fmt.Sprintf("%s/%s", project.Path, pAsset.Name))
		err = utils.Move(oldPath, newPath)

		if err != nil {
			log.Println("rename", err)
			return c.NoContent(http.StatusInternalServerError)
		}

		asset.Name = pAsset.Name
	}

	return c.NoContent(http.StatusOK)
}*/

func new(c echo.Context) error {

	pAsset := &models.ProjectAsset{}

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

	// Source
	src, err := files[0].Open()
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(fmt.Sprintf("%s/%s", path, files[0].Filename))
	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	asset, nestedAssets, err := models.NewProjectAsset(files[0].Filename, project, dst)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err = database.InsertAsset(asset); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	for _, a := range nestedAssets {
		if project.DefaultImagePath == "" && a.AssetType == "image" {
			project.DefaultImagePath = a.SHA1
			if err := database.UpdateProject(project); err != nil {
				log.Println(err)
			}
		}

		err := database.InsertAsset(a)
		if err != nil {
			log.Println(err)
		}
	}

	return c.JSON(http.StatusOK, asset)
}

func deleteAsset(c echo.Context) error {

	sha1 := c.Param("sha1")

	if sha1 == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := database.DeleteAsset(sha1); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}
