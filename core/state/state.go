package state

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/eduardooliveira/stLib/core/models"
	"github.com/eduardooliveira/stLib/core/runtime"
	"github.com/eduardooliveira/stLib/core/utils"
)

var TempFiles = make(map[string]*models.TempFile)
var Printers = make(map[string]*models.Printer)

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

func PercistPrinters() error {
	dataPath := runtime.Cfg.DataPath
	f, err := os.OpenFile(path.Join(dataPath, "printers.toml"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
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
	var err error

	dataPath := runtime.Cfg.DataPath
	if _, err = os.Stat(dataPath); os.IsNotExist(err) {
		err = utils.CreateFolder(dataPath)
		if err != nil {
			return err
		}
	}

	_, err = os.Stat(path.Join(dataPath, "printers.toml"))

	if err != nil {
		if _, err = os.Create(path.Join(dataPath, "printers.toml")); err != nil {
			return err
		}
	}

	_, err = toml.DecodeFile(path.Join(dataPath, "printers.toml"), &Printers)
	if err != nil {
		log.Println("error loading printers")
	}
	return err
}
