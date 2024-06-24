package data

import (
	"github.com/eduardooliveira/stLib/v2/database"
	"github.com/eduardooliveira/stLib/v2/library/entities"
)

func SaveAsset(a entities.Asset) error {
	return database.DB.Omit("NestedAssets").Save(&a).Error
}

func GetAsset(id string, deep bool) (rtn *entities.Asset, err error) {
	q := database.DB.Where(&entities.Asset{ID: id})
	if deep {
		q = q.Preload("NestedAssets.NestedAssets")
	}
	return rtn, q.First(&rtn).Error
}

func GetAssetByRootAndPath(root, path string, deep bool) (rtn *entities.Asset, err error) {
	q := database.DB.Where(&entities.Asset{Root: &root, Path: &path})
	if deep {
		q = q.Preload("NestedAssets.NestedAssets")
	}
	return rtn, q.First(&rtn).Error
}
