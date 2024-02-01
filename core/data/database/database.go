package database

import (
	"os"
	"path"

	"github.com/eduardooliveira/stLib/core/runtime"
	"github.com/eduardooliveira/stLib/core/utils"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() error {
	var err error

	dataPath := runtime.Cfg.DataPath
	if _, err = os.Stat(dataPath); os.IsNotExist(err) {
		err = utils.CreateFolder(dataPath)
		if err != nil {
			return err
		}
	}

	DB, err = gorm.Open(sqlite.Open(path.Join(dataPath, "data.db")), &gorm.Config{})
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

	return nil
}
