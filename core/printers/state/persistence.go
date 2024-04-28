package state

import (
	"log"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/eduardooliveira/stLib/core/printers/clients/klipper"
	"github.com/eduardooliveira/stLib/core/printers/entities"
	"github.com/eduardooliveira/stLib/core/runtime"
)

var Printers = make(map[string]entities.Printer)
var printerConfigs = make(map[string]*entities.Config)
var printersFile string

func Persist() error {
	f, err := os.OpenFile(printersFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
	if err := toml.NewEncoder(f).Encode(printerConfigs); err != nil {
		log.Println(err)
	}
	if err := f.Close(); err != nil {
		log.Println(err)
	}
	return err
}

func Load() error {
	printersFile = path.Join(runtime.GetDataPath(), "printers.toml")

	_, err := os.Stat(printersFile)

	if err != nil {
		if _, err = os.Create(printersFile); err != nil {
			return err
		}
	}

	_, err = toml.DecodeFile(printersFile, &printerConfigs)
	if err != nil {
		log.Println("error loading printers")
	}
	for uuid, config := range printerConfigs {
		if config.Type == "klipper" {
			Printers[uuid] = klipper.NewPrinter(config)
		}
	}
	return err
}
