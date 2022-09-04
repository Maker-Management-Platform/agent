package images

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/eduardooliveira/stLib/core/state"
	"github.com/labstack/echo/v4"
)

func get(c echo.Context) error {
	image, ok := state.Images[c.Param("sha1")]

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	return c.Attachment(image.Path, image.Name)
}

func upload(c echo.Context) error {
	projectUUID := c.FormValue("project")

	project, ok := state.Projects[projectUUID]
	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	file, err := c.FormFile("file")
	if err != nil {
		log.Println("Error getting the file: ", err)
		return err
	}
	src, err := file.Open()
	if err != nil {
		log.Println("Error opening the file: ", err)
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(fmt.Sprintf("%s/%s", project.Path, file.Filename))
	if err != nil {
		log.Println("Error creating the file: ", err)
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		log.Println("Error copying the file: ", err)
		return err
	}

	_, err = HandleImage(project, file.Filename)
	if err != nil {
		log.Println("Error handling the model: ", err)
		return err
	}

	err = state.PersistProject(project)
	if err != nil {
		log.Println("Error persisting the project: ", err)
		return err
	}

	return c.NoContent(http.StatusOK)
}
