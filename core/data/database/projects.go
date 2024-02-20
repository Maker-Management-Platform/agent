package database

import (
	models "github.com/eduardooliveira/stLib/core/entities"
	"gorm.io/gorm"
)

func initProjects() error {
	if err := DB.AutoMigrate(&models.Project{}); err != nil {
		return err
	}

	return DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Project{}).Error
}

func InsertProject(p *models.Project) error {
	return DB.Create(p).Error
}

func UpdateProject(p *models.Project) error {
	if err := DB.Save(p).Error; err != nil {
		return err
	}
	return DB.Model(p).Association("Tags").Replace(p.Tags)
}

func GetProjects() (rtn []*models.Project, err error) {
	return rtn, DB.Order("name").Find(&rtn).Error
}

func GetProjectNames() (rtn []*models.Project, err error) {
	return rtn, DB.Order("name").Select("uuid", "name").Find(&rtn).Error
}

func GetProject(uuid string) (rtn *models.Project, err error) {
	return rtn, DB.Where(&models.Project{UUID: uuid}).Preload("Tags").First(&rtn).Error
}

func DeleteProject(uuid string) (err error) {
	return DB.Where(&models.Project{UUID: uuid}).Delete(&models.Project{}).Error
}

func SetProjectDefaultImage(uuid string, imageId string) (err error) {
	return DB.Model(&models.Project{UUID: uuid}).Update("default_image_id", imageId).Error
}
