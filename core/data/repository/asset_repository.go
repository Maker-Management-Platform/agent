package repository

import (
	"context"

	"github.com/eduardooliveira/stLib/core/entities"
	"gorm.io/gorm"
)

type assetRepository struct {
	db *gorm.DB
}

// NewAssetRepository creates a new asset repository
func NewAssetRepository(db *gorm.DB) AssetRepository {
	return &assetRepository{db: db}
}

func (r *assetRepository) FindByID(ctx context.Context, id uint) (*entities.Asset, error) {
	var asset entities.Asset
	err := r.db.WithContext(ctx).
		Preload("Project").
		Preload("AssetType").
		First(&asset, id).Error

	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (r *assetRepository) FindAll(ctx context.Context) ([]*entities.Asset, error) {
	var assets []*entities.Asset
	err := r.db.WithContext(ctx).
		Preload("Project").
		Preload("AssetType").
		Find(&assets).Error

	if err != nil {
		return nil, err
	}

	return assets, nil
}

func (r *assetRepository) FindByProjectID(ctx context.Context, projectID uint) ([]*entities.Asset, error) {
	var assets []*entities.Asset
	err := r.db.WithContext(ctx).
		Preload("AssetType").
		Where("project_id = ?", projectID).
		Find(&assets).Error

	if err != nil {
		return nil, err
	}

	return assets, nil
}

func (r *assetRepository) Create(ctx context.Context, asset *entities.Asset) error {
	return r.db.WithContext(ctx).Create(asset).Error
}

func (r *assetRepository) Update(ctx context.Context, asset *entities.Asset) error {
	return r.db.WithContext(ctx).Save(asset).Error
}

func (r *assetRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entities.Asset{}, id).Error
}
