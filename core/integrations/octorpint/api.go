package octorpint

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eduardooliveira/stLib/core/data/database"
	"github.com/eduardooliveira/stLib/core/models"
	"github.com/eduardooliveira/stLib/core/utils"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func (p *OctoPrintPrinter) serverInfo() (*OctoPrintResponse, error) {
	bearer := "Bearer " + p.ApiKey
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/server", p.Address), nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	//TODO add error if forbidden
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var r OctoPrintResponse
	err = decoder.Decode(&r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (p *OctoPrintPrinter) ServerFilesUpload(asset *models.ProjectAsset) error {
	project, err := database.GetProject(asset.ProjectUUID)

	if err != nil {
		log.Println(err)
		return err
	}

	file, err := os.Open(utils.ToLibPath(fmt.Sprintf("%s/%s", project.FullPath(), asset.Name)))
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", asset.Name)
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	bearer := "Bearer " + p.ApiKey
	// location could be local or sdcard, can also create new folders by adding path
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/files/local", p.Address), body)

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if err != nil {
		return err
	} else {
		if resp.StatusCode != 201 {
			body := &bytes.Buffer{}
			body.ReadFrom(resp.Body)
			resp.Body.Close()
			fmt.Println(resp.StatusCode)
			fmt.Println(resp.Header)
			fmt.Println(body)
			return errors.New("unknown error uploading file")
		}
	}

	return nil
}
