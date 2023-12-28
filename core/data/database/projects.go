package database

import (
	"github.com/eduardooliveira/stLib/core/models"
	"gorm.io/gorm"
)

func initProjects() error {
	if err := db.AutoMigrate(&models.Project{}); err != nil {
		return err
	}

	return db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Project{}).Error
}

func InsertProject(p *models.Project) error {
	return db.Create(p).Error
}

func UpdateProject(p *models.Project) error {
	return db.Save(p).Error
}

func GetProjects() (rtn []*models.Project, err error) {
	return rtn, db.Order("name").Find(&rtn).Error
}

func GetProject(uuid string) (rtn *models.Project, err error) {
	return rtn, db.Where(&models.Project{UUID: uuid}).First(&rtn).Error
}
