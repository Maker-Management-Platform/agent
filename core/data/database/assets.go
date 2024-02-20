package database

import (
	models "github.com/eduardooliveira/stLib/core/entities"
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

func GetAsset(id string) (rtn *models.ProjectAsset, err error) {
	return rtn, DB.Where(&models.ProjectAsset{ID: id}).First(&rtn).Error
}

func GetProjectAsset(uuid string, id string) (rtn *models.ProjectAsset, err error) {
	return rtn, DB.Where(&models.ProjectAsset{ID: id, ProjectUUID: uuid}).First(&rtn).Error
}

func DeleteAsset(id string) (err error) {
	return DB.Where(&models.ProjectAsset{ID: id}).Delete(&models.ProjectAsset{}).Error
}

func SetModelImage(id string, imageId string) (err error) {
	return DB.Model(&models.ProjectAsset{ID: id}).Update("model.image_id", imageId).Error
}
