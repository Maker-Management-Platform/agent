package projects

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/eduardooliveira/stLib/core/entities"
	"github.com/eduardooliveira/stLib/core/utils"
	"github.com/labstack/echo/v4"
)

type CreateProject struct {
	Name             string          `json:"name"`
	Description      string          `json:"description"`
	DefaultImageName string          `json:"default_image_name"`
	Tags             []*entities.Tag `json:"tags"`
}

func new(c echo.Context) error {

	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	files := form.File["files"]

	if len(files) == 0 {
		log.Println("No files")
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("no files uploaded").Error())
	}

	projectPayload := form.Value["payload"]
	if len(projectPayload) != 1 {
		log.Println("more payloads than expected")
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("more payloads than expected").Error())
	}

	createProject := &CreateProject{}
	err = json.Unmarshal([]byte(projectPayload[0]), createProject)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	projectFolder := filepath.Clean(createProject.Name)

	path := utils.ToLibPath(projectFolder)
	if err := utils.CreateFolder(path); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	for _, file := range files {
		// Source
		src, err := file.Open()
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		defer src.Close()

		// Destination
		dst, err := os.Create(filepath.Join(path, file.Filename))
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

	}

	/*var project *entities.Project
	if project, err = processing.HandlePath(projectFolder); err != nil {
		log.Printf("error loading the project %q: %v\n", path, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	project.Description = createProject.Description
	project.Tags = createProject.Tags

	if err = database.UpdateProject(project); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, struct {
		UUID string `json:"uuid"`
	}{project.UUID})*/
	return c.JSON(http.StatusOK, struct {
		UUID string `json:"uuid"`
	}{"wqe"})
}
