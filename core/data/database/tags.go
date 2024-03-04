package database

import (
	models "github.com/eduardooliveira/stLib/core/entities"
)

func initTags() error {
	return DB.AutoMigrate(&models.Tag{})
}

func GetTags() (rtn []*models.Tag, err error) {
	return rtn, DB.Order("value").Find(&rtn).Error
}
