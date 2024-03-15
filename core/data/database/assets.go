package database

import (
	"errors"

	"github.com/eduardooliveira/stLib/core/entities"
	"github.com/eduardooliveira/stLib/core/system"
	"gorm.io/gorm"
)

const assetEvent = "asset.event"

func initAssets() error {
	return DB.AutoMigrate(&entities.ProjectAsset{})
}

func InsertAsset(a *entities.ProjectAsset) error {
	if err := DB.Create(a).Error; err != nil {
		return err

	}
	system.Publish(assetEvent, map[string]any{"projectUUID": a.ProjectUUID, "assetID": a.ID, "assetLabel": a.Label, "assetName": a.Name, "type": "new"})
	return nil
}

func SaveAsset(a *entities.ProjectAsset) error {
	if err := DB.Create(a).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			if err := DB.Save(a).Error; err != nil {
				return err
			} else {
				system.Publish(assetEvent, map[string]any{"projectUUID": a.ProjectUUID, "assetID": a.ID, "assetLabel": a.Label, "assetName": a.Name, "type": "update"})
			}
		} else {
			return err
		}
	}
	system.Publish(assetEvent, map[string]any{"projectUUID": a.ProjectUUID, "assetID": a.ID, "assetLabel": a.Label, "assetName": a.Name, "type": "new"})

	return nil
}

func GetAssetsByProject(uuid string) (rtn []*entities.ProjectAsset, err error) {
	return rtn, DB.Order("name").Where(&entities.ProjectAsset{ProjectUUID: uuid}).Find(&rtn).Error
}

func GetAsset(id string) (rtn *entities.ProjectAsset, err error) {
	return rtn, DB.Where(&entities.ProjectAsset{ID: id}).First(&rtn).Error
}

func GetAssetByProjectAndName(uuid string, name string) (rtn *entities.ProjectAsset, err error) {
	return rtn, DB.Debug().Where(&entities.ProjectAsset{ProjectUUID: uuid, Name: name}).First(&rtn).Error
}

func GetProjectAsset(uuid string, id string) (rtn *entities.ProjectAsset, err error) {
	return rtn, DB.Where(&entities.ProjectAsset{ID: id, ProjectUUID: uuid}).First(&rtn).Error
}

func DeleteAsset(id string) (err error) {
	return DB.Where(&entities.ProjectAsset{ID: id}).Delete(&entities.ProjectAsset{}).Error
}

func UpdateAssetImage(a *entities.ProjectAsset, imageID string) (err error) {
	if err := DB.Model(&entities.ProjectAsset{ID: a.ID}).Update("image_id", imageID).Error; err != nil {
		return err
	}
	system.Publish(assetEvent, map[string]any{"projectUUID": a.ProjectUUID, "assetID": a.ID, "assetLabel": a.Label, "assetName": a.Name, "type": "update"})
	return nil
}

func UpdateAssetProperties(a *entities.ProjectAsset, properties entities.AssetProperties) error {
	if err := DB.Model(&entities.ProjectAsset{ID: a.ID}).Update("properties", properties).Error; err != nil {
		return err
	}
	system.Publish(assetEvent, map[string]any{"projectUUID": a.ProjectUUID, "assetID": a.ID, "assetLabel": a.Label, "assetName": a.Name, "type": "update"})
	return nil
}
