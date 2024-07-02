package repo

import (
	"errors"
	"math"

	"github.com/eduardooliveira/stLib/v2/database"
	"github.com/eduardooliveira/stLib/v2/library/entities"
	"gorm.io/gorm"
)

func (r AssetRepo) SaveAsset(a entities.Asset) error {
	return database.DB.Omit("NestedAssets").Save(&a).Error
}

func (r AssetRepo) GetAsset(id string, deep bool) (rtn *entities.Asset, err error) {
	q := database.DB.Where(&entities.Asset{ID: id})
	if deep {
		q = q.Preload("NestedAssets.NestedAssets")
	}
	return rtn, q.First(&rtn).Error
}

func (r AssetRepo) GetAssetByRootAndPath(root, path string, deep bool) (rtn *entities.Asset, err error) {
	q := database.DB.Where(&entities.Asset{Root: &root, Path: &path})
	if deep {
		q = q.Preload("NestedAssets.NestedAssets")
	}
	return rtn, q.First(&rtn).Error
}

func (r AssetRepo) LoadParents(a *entities.Asset, dept int, fields ...string) error {
	if dept == 0 {
		return errors.New("dept must be greater than 0")
	}
	q := database.DB.Model(a)
	for i := 0; i < dept; i++ {
		q = q.Preload("Parent", func(q *gorm.DB) *gorm.DB {
			if len(fields) > 0 {
				return q.Select(fields)
			}
			return q
		})
	}
	return q.Find(a).Error
}

type GetAssetParams struct {
	Root    string
	Path    string
	PerPage int
	Page    int
}

func (r AssetRepo) GetPagedNested(asset *entities.Asset, page, perPage int) (int, error) {
	var totalRows int64
	if err := database.DB.Model(entities.Asset{}).
		Where("parent_id = ?", asset.ID).
		Count(&totalRows).Error; err != nil {
		return 0, err
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(perPage)))

	return totalPages, database.DB.
		Where("parent_id = ?", asset.ID).
		Offset(page * perPage).
		Limit(perPage).
		Order("Label ASC").
		Find(&asset.NestedAssets).Error
}
