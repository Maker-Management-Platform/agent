package state

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/eduardooliveira/stLib/core/models"
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
	f, err := os.OpenFile("data/printers.toml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
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

	err := utils.CreateFolder("data")

	if err != nil {
		return err
	}

	_, err = os.Stat("data/printers.toml")

	if err != nil {
		if _, err = os.Create("data/printers.toml"); err != nil {
			return err
		}
	}

	_, err = toml.DecodeFile("data/printers.toml", &Printers)
	if err != nil {
		log.Println("error loading printers")
	}
	return err
}
