package state

import (
	"log"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/eduardooliveira/stLib/core/entities"
	"github.com/eduardooliveira/stLib/core/runtime"
)

var TempFiles = make(map[string]*entities.TempFile)
var Printers = make(map[string]*entities.Printer)
var AssetTypes = make(map[string]*entities.AssetType)
var ExtensionProjectType = make(map[string]*entities.AssetType)
var printersFile string
var assetTypesFile string

func PersistPrinters() error {
	f, err := os.OpenFile(printersFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
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

func LoadAssetTypes() error {
	assetTypesFile = path.Join(runtime.GetDataPath(), "assetTypes.toml")

	_, err := os.Stat(assetTypesFile)

	if err != nil {
		if _, err = os.Create(assetTypesFile); err != nil {
			return err
		}
	}

	_, err = toml.DecodeFile(assetTypesFile, &AssetTypes)
	if err != nil {
		log.Println("error loading asset types")
	}

	if len(AssetTypes) == 0 {
		AssetTypes["model"] = &entities.AssetType{
			Name:       "model",
			Label:      "Models",
			Extensions: []string{".stl", ".3fm"},
			Order:      0,
		}
		AssetTypes["image"] = &entities.AssetType{
			Name:       "image",
			Label:      "Images",
			Extensions: []string{".png", ".jpg", ".jpeg", ".gif", ".bmp", ".webp"},
			Order:      1,
		}
		AssetTypes["slice"] = &entities.AssetType{
			Name:       "slice",
			Label:      "Slices",
			Extensions: []string{".gcode"},
			Order:      2,
		}
		AssetTypes["source"] = &entities.AssetType{
			Name:       "source",
			Label:      "Sources",
			Extensions: []string{".stp", ".step", ".ste", ".fbx", ".f3d", ".f3z", ".iam", ".ipt"},
			Order:      99,
		}
		f, err := os.OpenFile(printersFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
		if err := toml.NewEncoder(f).Encode(Printers); err != nil {
			log.Println(err)
		}
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}

	for _, assetType := range AssetTypes {
		for _, ext := range assetType.Extensions {
			ExtensionProjectType[ext] = assetType
		}
	}
	return err
}
