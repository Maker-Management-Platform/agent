package database

import (
	"path"

	"github.com/eduardooliveira/stLib/core/runtime"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() error {
	var err error

	DB, err = gorm.Open(sqlite.Open(path.Join(runtime.GetDataPath(), "data.db")), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		panic("failed to connect database")
	}

	if err = initTags(); err != nil {
		return err
	}

	if err = initProjects(); err != nil {
		return err
	}

	if err = initAssets(); err != nil {
		return err
	}

	if err = initPrintJob(); err != nil {
		return err
	}

	return nil
}
