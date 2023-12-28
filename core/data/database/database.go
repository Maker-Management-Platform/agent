package database

import (
	"github.com/eduardooliveira/stLib/core/utils"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDatabase() error {
	err := utils.CreateFolder("data")
	if err != nil {
		return err
	}

	db, err = gorm.Open(sqlite.Open("data/data.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if err = initProjects(); err != nil {
		return err
	}

	if err = initAssets(); err != nil {
		return err
	}

	return nil
}
