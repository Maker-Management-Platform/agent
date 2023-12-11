package slicer

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/eduardooliveira/stLib/core/models"
	"github.com/eduardooliveira/stLib/core/state"
	"github.com/labstack/echo/v4"
)

func version(c echo.Context) error {
	return c.JSON(http.StatusOK, struct {
		API    string `json:"api"`
		Server string `json:"server"`
		Text   string `json:"text"`
	}{
		API:    "xxx",
		Server: "xx",
		Text:   "OctoPrint xx",
	})
}

func upload(c echo.Context) error {

	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}

	files := form.File["file"]

	if len(files) == 0 {
		log.Println("No files")
		return c.NoContent(http.StatusBadRequest)
	}
	name := files[0].Filename

	// Source
	src, err := files[0].Open()
	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(fmt.Sprintf("%s/%s", "temp", name))
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

	tempFile, _ := models.NewTempFile(name)

	token := strings.Split(strings.ToLower(name), "_")[0]

	for _, p := range state.Projects {
		if strings.Contains(strings.ToLower(p.Name), token) {
			tempFile.AddMatch(p.UUID)
			log.Println(fmt.Sprintf("Found project match: %s -> %s", p.Name, token))
		}
		for _, tag := range p.Tags {
			if strings.Contains(strings.ToLower(tag), token) {
				tempFile.AddMatch(p.UUID)
				log.Println(fmt.Sprintf("Found tag match: %s -> %s", tag, token))
			}
		}
	}

	state.TempFiles[tempFile.UUID] = tempFile

	return c.NoContent(http.StatusOK)
}
