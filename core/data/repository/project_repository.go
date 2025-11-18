package repository

import (
	"context"

	"github.com/eduardooliveira/stLib/core/entities"
	"gorm.io/gorm"
)

type projectRepository struct {
	db *gorm.DB
}

// NewProjectRepository creates a new project repository
func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) FindByID(ctx context.Context, id uint) (*entities.Project, error) {
	var project entities.Project
	err := r.db.WithContext(ctx).
		Preload("Assets").
		Preload("Tags").
		First(&project, id).Error

	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (r *projectRepository) FindAll(ctx context.Context) ([]*entities.Project, error) {
	var projects []*entities.Project
	err := r.db.WithContext(ctx).
		Preload("Assets").
		Preload("Tags").
		Find(&projects).Error

	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *projectRepository) Create(ctx context.Context, project *entities.Project) error {
	return r.db.WithContext(ctx).Create(project).Error
}

func (r *projectRepository) Update(ctx context.Context, project *entities.Project) error {
	return r.db.WithContext(ctx).Save(project).Error
}

func (r *projectRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entities.Project{}, id).Error
}

func (r *projectRepository) FindByPath(ctx context.Context, path string) (*entities.Project, error) {
	var project entities.Project
	err := r.db.WithContext(ctx).
		Preload("Assets").
		Preload("Tags").
		Where("path = ?", path).
		First(&project).Error

	if err != nil {
		return nil, err
	}

	return &project, nil
}
