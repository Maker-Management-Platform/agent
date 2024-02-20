package processing

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/eduardooliveira/stLib/core/data/database"
	models "github.com/eduardooliveira/stLib/core/entities"
	"github.com/eduardooliveira/stLib/core/runtime"
	"github.com/eduardooliveira/stLib/core/state"
)

func RunTempDiscovery() {
	log.Println("Discovering Temp files")

	tempPath := filepath.Clean(path.Join(runtime.GetDataPath(), "temp"))
	if _, err := os.Stat(tempPath); os.IsNotExist(err) {
		err := os.MkdirAll(tempPath, os.ModePerm)
		if err != nil {
			log.Panic(err)
		}
	}

	entries, err := os.ReadDir(tempPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		blacklisted := false
		for _, blacklist := range runtime.Cfg.Library.Blacklist {
			if strings.HasSuffix(e.Name(), blacklist) {
				blacklisted = true
				break
			}
		}
		if blacklisted {
			continue
		}
		fmt.Println(e.Name())
		tempFile, err := DiscoverTempFile(e.Name())
		if err != nil {
			log.Println("Error Discovering temp file: ", err)
			continue
		}
		state.TempFiles[tempFile.UUID] = tempFile
	}
}

func DiscoverTempFile(name string) (*models.TempFile, error) {
	tempFile, err := models.NewTempFile(name)
	if err != nil {
		return nil, err
	}

	token := strings.Split(strings.ToLower(name), "_")[0]

	projects, err := database.GetProjects()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for _, p := range projects {
		if strings.Contains(strings.ToLower(p.Name), token) {
			tempFile.AddMatch(p.UUID)
		}
		for _, tag := range p.Tags {
			if strings.Contains(strings.ToLower(tag.Value), token) {
				tempFile.AddMatch(p.UUID)
			}
		}
	}
	return tempFile, nil
}
