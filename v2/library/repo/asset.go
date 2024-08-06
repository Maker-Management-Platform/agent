package repo

import (
	"errors"
	"math"

	"github.com/eduardooliveira/stLib/v2/database"
	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/utils"
	"gorm.io/gorm"
)

func (r AssetRepo) SaveAsset(a entities.Asset) error {
	return database.DB.Omit("NestedAssets").Save(&a).Error
}

func (r AssetRepo) GetAsset(id string, deep bool) (rtn entities.Asset, err error) {
	q := database.DB.Where("ID", id)
	if deep {
		q = q.Preload("NestedAssets.NestedAssets")
	}
	return rtn, q.First(&rtn).Error
}

func (r AssetRepo) GetAssetRoots(deep bool) (rtn []*entities.Asset, err error) {
	q := database.DB.Debug().Where(&entities.Asset{NodeKind: utils.Ptr(entities.NodeKindRoot)})
	if deep {
		q = q.Preload("NestedAssets.NestedAssets")
	}
	return rtn, q.Find(&rtn).Error
}

func (r AssetRepo) GetAssetByRootAndPath(root, path string, deep bool) (rtn entities.Asset, err error) {
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

	parentChain := "Parent"
	for i := 0; i < dept; i++ {
		parentChain += ".Parent"
	}
	q = q.Preload(parentChain, func(q *gorm.DB) *gorm.DB {
		if len(fields) > 0 {
			fields = append(fields, "ParentID")
			return q.Select(fields)
		}
		return q
	})
	return q.Debug().Find(a).Error
}

func (r AssetRepo) GetPagedNested(asset, filter *entities.Asset, page, perPage int) (int, error) {
	var totalRows int64
	if err := database.DB.Model(entities.Asset{}).
		Where(filter).
		Count(&totalRows).Error; err != nil {
		return 0, err
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(perPage)))

	return totalPages, database.DB.
		Where(filter).
		Offset(page * perPage).
		Limit(perPage).
		Order("Label ASC").
		Find(&asset.NestedAssets).Error
}

func (r AssetRepo) SetDirtyRoot(root string) error {
	return database.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).
		Model(&entities.Asset{}).Update("SeenOnScan", false).Error
}

func (r AssetRepo) SetDirtyFS(fsName string) error {
	return database.DB.Model(&entities.Asset{}).Where("fs_name", fsName).Update("SeenOnScan", false).Error
}

func (r AssetRepo) DeleteUnSeenInFS(fsName string) error {
	return database.DB.Model(entities.Asset{}).
		Delete(entities.Asset{}, entities.Asset{SeenOnScan: utils.Ptr(false), FSName: &fsName}).Error
}

func (r AssetRepo) DeleteUnSeenInRoot(root string) error {
	return database.DB.Model(entities.Asset{}).
		Delete(entities.Asset{}, entities.Asset{SeenOnScan: utils.Ptr(false), Root: &root}).Error
}

func (r AssetRepo) UpdateAsset(a *entities.Asset) error {
	return database.DB.Model(&entities.Asset{ID: a.ID}).Updates(a).Error
}
