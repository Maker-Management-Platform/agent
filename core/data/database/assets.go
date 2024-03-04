package database

import (
	"github.com/eduardooliveira/stLib/core/entities"
)

func initAssets() error {
	return DB.AutoMigrate(&entities.ProjectAsset{})
}

func InsertAsset(a *entities.ProjectAsset) error {
	return DB.Create(a).Error
}

func GetAssetsByProject(uuid string) (rtn []*entities.ProjectAsset, err error) {
	return rtn, DB.Order("name").Where(&entities.ProjectAsset{ProjectUUID: uuid}).Find(&rtn).Error
}

func GetAsset(id string) (rtn *entities.ProjectAsset, err error) {
	return rtn, DB.Where(&entities.ProjectAsset{ID: id}).First(&rtn).Error
}

func GetAssetByProjectAndName(uuid string, name string) (rtn *entities.ProjectAsset, err error) {
	return rtn, DB.Where(&entities.ProjectAsset{ProjectUUID: uuid, Name: name}).First(&rtn).Error
}

func GetProjectAsset(uuid string, id string) (rtn *entities.ProjectAsset, err error) {
	return rtn, DB.Where(&entities.ProjectAsset{ID: id, ProjectUUID: uuid}).First(&rtn).Error
}

func DeleteAsset(id string) (err error) {
	return DB.Where(&entities.ProjectAsset{ID: id}).Delete(&entities.ProjectAsset{}).Error
}

func UpdateAssetImage(id string, imageID string) (err error) {
	return DB.Model(&entities.ProjectAsset{ID: id}).Update("image_id", imageID).Error
}

func UpdateAssetProperties(id string, properties entities.AssetProperties) error {
	return DB.Model(&entities.ProjectAsset{ID: id}).Update("properties", properties).Error
}
