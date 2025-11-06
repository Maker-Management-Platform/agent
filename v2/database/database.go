package database

import (
	"log"
	"path"

	"github.com/eduardooliveira/stLib/v2/config"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() error {
	var err error
	DB, err = gorm.Open(sqlite.Open(path.Join(config.Cfg.Core.DataFolder, "data2.db")), &gorm.Config{
		TranslateError: true,
		Logger:         logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("failed to connect database")
	}
	return nil
}
