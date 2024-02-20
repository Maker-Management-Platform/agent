package state

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	models "github.com/eduardooliveira/stLib/core/entities"
	"github.com/eduardooliveira/stLib/core/runtime"
	"github.com/eduardooliveira/stLib/core/utils"
)

var TempFiles = make(map[string]*models.TempFile)
var Printers = make(map[string]*models.Printer)
var printersFile string

func PersistProject(project *models.Project) error {
	f, err := os.OpenFile(fmt.Sprintf("%s/.project.stlib", utils.ToLibPath(project.FullPath())), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Println(err)
	}
	if err := toml.NewEncoder(f).Encode(project); err != nil {
		log.Println(err)
	}
	if err := f.Close(); err != nil {
		log.Println(err)
	}
	return err
}

func PersistPrinters() error {
	f, err := os.OpenFile(printersFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Println(err)
	}
	if err := toml.NewEncoder(f).Encode(Printers); err != nil {
		log.Println(err)
	}
	if err := f.Close(); err != nil {
		log.Println(err)
	}
	return err
}

func LoadPrinters() error {
	printersFile = path.Join(runtime.GetDataPath(), "printers.toml")

	_, err := os.Stat(printersFile)

	if err != nil {
		if _, err = os.Create(printersFile); err != nil {
			return err
		}
	}

	_, err = toml.DecodeFile(printersFile, &Printers)
	if err != nil {
		log.Println("error loading printers")
	}
	return err
}
