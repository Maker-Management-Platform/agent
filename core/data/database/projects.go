package database

import (
	"github.com/eduardooliveira/stLib/core/models"
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
	return DB.Save(p).Error
}

func GetProjects() (rtn []*models.Project, err error) {
	return rtn, DB.Order("name").Find(&rtn).Error
}

func GetProject(uuid string) (rtn *models.Project, err error) {
	return rtn, DB.Where(&models.Project{UUID: uuid}).First(&rtn).Error
}
