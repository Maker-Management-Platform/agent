package klipper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/eduardooliveira/stLib/core/data/database"
	"github.com/eduardooliveira/stLib/core/entities"
	"github.com/eduardooliveira/stLib/core/utils"
)

func (p *KlipperPrinter) serverInfo() (*Result, error) {
	res, err := http.Get(fmt.Sprintf("%s/server/info", p.config.Address))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r MoonRakerResponse
	err = decoder.Decode(&r)
	if err != nil {
		return nil, err
	}

	return r.Result, nil
}

func (p *KlipperPrinter) serverFilesUpload(asset *entities.ProjectAsset) error {

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

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/server/files/upload", p.config.Address), body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
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
