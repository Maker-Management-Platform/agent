package assets

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/eduardooliveira/stLib/core/discovery"
	"github.com/eduardooliveira/stLib/core/models"
	"github.com/eduardooliveira/stLib/core/state"
	"github.com/eduardooliveira/stLib/core/utils"
	"github.com/labstack/echo/v4"
)

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
}

func new(c echo.Context) error {

	pAsset := &models.ProjectAsset{}
	err := c.Bind(pAsset)

	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	form, err := c.MultipartForm()

	files := form.File["files"]

	if len(files) == 0 {
		log.Println("No files")
		return c.NoContent(http.StatusBadRequest)
	}

	project, ok := state.Projects[pAsset.ProjectUUID]

	if !ok {
		return c.NoContent(http.StatusInternalServerError)
	}

	path := utils.ToLibPath(fmt.Sprintf("%s/%s", project.FullPath(), pAsset.Name))

	// Source
	src, err := files[0].Open()
	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
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

	err = discovery.DiscoverProjectAssets(project)
	if err != nil {
		log.Printf("error loading the project %q: %v\n", path, err)
		return err
	}

	for _, a := range state.Projects[project.UUID].Assets {
		if a.Name == files[0].Filename {
			return c.JSON(http.StatusOK, a)
		}
	}

	return c.NoContent(http.StatusInternalServerError)
}

func deleteAsset(c echo.Context) error {

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

	err := os.Remove(utils.ToLibPath(fmt.Sprintf("%s/%s", project.FullPath(), asset.Name)))
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	delete(state.Assets, sha1)
	delete(project.Assets, sha1)

	return c.NoContent(http.StatusOK)
}
