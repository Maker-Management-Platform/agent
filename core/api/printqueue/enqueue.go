package printqueue

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/eduardooliveira/stLib/core/data/database"
	"github.com/eduardooliveira/stLib/core/entities"
	"github.com/labstack/echo/v4"
)

type enqueueRequest struct {
	Instances   int      `json:"instances"`
	ProjectId   string   `json:"projectId"`
	SliceId     string   `json:"sliceId"`
	PrinterUUID string   `json:"printerUUID"`
	Tags        []string `json:"tags"`
}

func enqueue(c echo.Context) error {

	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	enqueuePayload := form.Value["payload"]
	if len(enqueuePayload) != 1 {
		log.Println("more payloads than expected")
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("more payloads than expected").Error())
	}

	enqueueRequest := &enqueueRequest{}
	err = json.Unmarshal([]byte(enqueuePayload[0]), enqueueRequest)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	files := form.File["files"]

	if enqueueRequest.SliceId == "" && len(files) == 0 {
		log.Println("No files")
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("no files uploaded").Error())
	}

	var slice *entities.ProjectAsset

	if enqueueRequest.SliceId != "" {
		slice, err = database.GetAsset(enqueueRequest.SliceId)
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	} else {
		log.Println("handle slice upload")
	}

	project, err := database.GetProject(slice.ProjectUUID)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	log.Println(project)

	for i := 0; i < enqueueRequest.Instances; i++ {
		job := entities.NewPrintJob(slice)
		err := database.InsertPrintJob(job)
		if err != nil {
			log.Println(err)
		}
	}

	return c.NoContent(http.StatusCreated)
}
