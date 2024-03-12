package database

import (
	"github.com/eduardooliveira/stLib/core/entities"
	"github.com/eduardooliveira/stLib/core/system"
)

const projectEvent = "project.event"

func initProjects() error {
	return DB.AutoMigrate(&entities.Project{})
}

func InsertProject(p *entities.Project) error {
	if err := DB.Create(p).Error; err != nil {
		return err
	}
	system.Publish(projectEvent, map[string]any{"projectUUID": p.UUID, "projectName": p.Name, "type": "new"})
	return nil
}

func UpdateProject(p *entities.Project) error {
	if err := DB.Save(p).Error; err != nil {
		return err
	}
	if err := DB.Model(p).Association("Tags").Replace(p.Tags); err != nil {
		return err
	}
	system.Publish(projectEvent, map[string]any{"projectUUID": p.UUID, "projectName": p.Name, "type": "update"})
	return nil
}

func GetProjects() (rtn []*entities.Project, err error) {
	return rtn, DB.Order("name").Find(&rtn).Error
}

func GetProjectNames() (rtn []*entities.Project, err error) {
	return rtn, DB.Order("name").Select("uuid", "name").Find(&rtn).Error
}

func GetProject(uuid string) (rtn *entities.Project, err error) {
	return rtn, DB.Where(&entities.Project{UUID: uuid}).Preload("Tags").First(&rtn).Error
}

func GetProjectByPathAndName(path string, name string) (rtn *entities.Project, err error) {
	return rtn, DB.Where(&entities.Project{Path: path, Name: name}).First(&rtn).Error
}

func DeleteProject(uuid string) (err error) {
	return DB.Where(&entities.Project{UUID: uuid}).Delete(&entities.Project{}).Error
}

func SetProjectDefaultImage(p *entities.Project, imageId string) (err error) {
	if err := DB.Model(&entities.Project{UUID: p.UUID}).Update("default_image_id", imageId).Error; err != nil {
		return err
	}
	system.Publish(projectEvent, map[string]any{"projectUUID": p.UUID, "projectName": p.Name, "type": "update"})
	return nil
}
