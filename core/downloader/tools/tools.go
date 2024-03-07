package tools

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	models "github.com/eduardooliveira/stLib/core/entities"
	"github.com/eduardooliveira/stLib/core/processing"
	"github.com/eduardooliveira/stLib/core/utils"
)

func DownloadAsset(name string, project *models.Project, client *http.Client, req *http.Request) error {
	out, err := os.Create(utils.ToLibPath(fmt.Sprintf("%s/%s", project.FullPath(), name)))
	if err != nil {
		return err
	}
	defer out.Close()

	log.Println("Downloading: ", name)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	processing.EnqueueInitJob(&processing.ProcessableAsset{
		Name:    name,
		Project: project,
		Origin:  "fs",
	})

	return nil
}

func SaveFile(dst string, f io.Reader) error {
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, f)
	if err != nil {
		return err
	}

	return nil
}
