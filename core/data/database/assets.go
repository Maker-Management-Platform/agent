package database

import (
	"github.com/eduardooliveira/stLib/core/models"
	"gorm.io/gorm"
)

func initAssets() error {
	if err := DB.AutoMigrate(&models.ProjectAsset{}); err != nil {
		return err
	}

	return DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.ProjectAsset{}).Error
}

func InsertAsset(a *models.ProjectAsset) error {
	return DB.Create(a).Error
}

func GetAssetsByProject(uuid string) (rtn []*models.ProjectAsset, err error) {
	return rtn, DB.Order("name").Where(&models.ProjectAsset{ProjectUUID: uuid}).Find(&rtn).Error
}

func GetAsset(uuid string, sha1 string) (rtn *models.ProjectAsset, err error) {
	return rtn, DB.Where(&models.ProjectAsset{SHA1: sha1, ProjectUUID: uuid}).First(&rtn).Error
}
