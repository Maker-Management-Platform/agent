package discovery

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/eduardooliveira/stLib/core/models"
	"github.com/eduardooliveira/stLib/core/runtime"
	"github.com/eduardooliveira/stLib/core/state"
)

func RunTempDiscovery() {
	log.Println("Discovering Temp files")
	entries, err := os.ReadDir("temp")
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		blacklisted := false
		for _, blacklist := range runtime.Cfg.FileBlacklist {
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

	for _, p := range state.Projects {
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
