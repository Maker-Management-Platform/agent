package database

import (
	"github.com/eduardooliveira/stLib/core/models"
	"gorm.io/gorm"
)

func initTags() error {
	if err := DB.AutoMigrate(&models.Tag{}); err != nil {
		return err
	}

	return DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Tag{}).Error
}

func GetTags() (rtn []*models.Tag, err error) {
	return rtn, DB.Order("value").Find(&rtn).Error
}
