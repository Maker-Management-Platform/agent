package repository

import (
	"context"

	"github.com/eduardooliveira/stLib/core/entities"
	"gorm.io/gorm"
)

type tagRepository struct {
	db *gorm.DB
}

// NewTagRepository creates a new tag repository
func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) FindByID(ctx context.Context, id uint) (*entities.Tag, error) {
	var tag entities.Tag
	err := r.db.WithContext(ctx).
		Preload("Projects").
		First(&tag, id).Error

	if err != nil {
		return nil, err
	}

	return &tag, nil
}

func (r *tagRepository) FindAll(ctx context.Context) ([]*entities.Tag, error) {
	var tags []*entities.Tag
	err := r.db.WithContext(ctx).Find(&tags).Error

	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (r *tagRepository) FindByName(ctx context.Context, name string) (*entities.Tag, error) {
	var tag entities.Tag
	err := r.db.WithContext(ctx).
		Where("name = ?", name).
		First(&tag).Error

	if err != nil {
		return nil, err
	}

	return &tag, nil
}

func (r *tagRepository) Create(ctx context.Context, tag *entities.Tag) error {
	return r.db.WithContext(ctx).Create(tag).Error
}

func (r *tagRepository) Update(ctx context.Context, tag *entities.Tag) error {
	return r.db.WithContext(ctx).Save(tag).Error
}

func (r *tagRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entities.Tag{}, id).Error
}
