package integration

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"path/filepath"

	v1entities "github.com/eduardooliveira/stLib/core/entities"
	v2entities "github.com/eduardooliveira/stLib/v2/library/entities"
	"gorm.io/gorm"
)

// DataMigration handles migration from v1 to v2 data models
type DataMigration struct {
	v1DB *gorm.DB
	v2DB *gorm.DB
}

// NewDataMigration creates a new data migration instance
func NewDataMigration(v1DB, v2DB *gorm.DB) *DataMigration {
	return &DataMigration{
		v1DB: v1DB,
		v2DB: v2DB,
	}
}

// MigrateProject converts a v1 Project to v2 Asset (root asset)
func (m *DataMigration) MigrateProject(v1Project *v1entities.Project, fs v2entities.LibFS) (*v2entities.Asset, error) {
	// Generate ID from path
	data := []byte(filepath.Join(fs.GetName(), fs.GetRoot(), v1Project.Path))
	md5Hash := md5.Sum(data)
	id := hex.EncodeToString(md5Hash[:])

	// Create root asset
	asset := &v2entities.Asset{
		ID:       id,
		Label:    &v1Project.Name,
		Path:     &v1Project.Path,
		Root:     fs.GetRoot(),
		FSKind:   fs.Kind(),
		FSName:   fs.GetName(),
		NodeKind: v2entities.NodeKindRoot,
	}

	// Set description if available
	if v1Project.Description != "" {
		asset.Description = &v1Project.Description
	}

	// Migrate tags
	for _, v1Tag := range v1Project.Tags {
		v2Tag := &v2entities.Tag{
			Name: v1Tag.Name,
		}
		asset.Tags = append(asset.Tags, v2Tag)
	}

	return asset, nil
}

// MigrateAsset converts a v1 Asset to v2 Asset (nested asset)
func (m *DataMigration) MigrateAsset(v1Asset *v1entities.Asset, fs v2entities.LibFS, parentID *string) (*v2entities.Asset, error) {
	// Generate ID from path
	data := []byte(filepath.Join(fs.GetName(), fs.GetRoot(), v1Asset.Path))
	md5Hash := md5.Sum(data)
	id := hex.EncodeToString(md5Hash[:])

	ext := filepath.Ext(v1Asset.Path)

	asset := &v2entities.Asset{
		ID:        id,
		Label:     &v1Asset.Name,
		Path:      &v1Asset.Path,
		Root:      fs.GetRoot(),
		FSKind:    fs.Kind(),
		FSName:    fs.GetName(),
		Extension: &ext,
		NodeKind:  v2entities.NodeKindFile,
		ParentID:  parentID,
	}

	// Set kind based on asset type
	if v1Asset.AssetType != nil {
		asset.Kind = &v1Asset.AssetType.Name
	}

	// Set thumbnail if available
	if v1Asset.Thumbnail != "" {
		asset.Thumbnail = &v1Asset.Thumbnail
	}

	// Migrate properties from v1 asset fields to v2 properties
	properties := make(v2entities.Properties)

	if v1Asset.RenderedImage != "" {
		properties["renderedImage"] = v1Asset.RenderedImage
	}

	if v1Asset.FileSize > 0 {
		properties["fileSize"] = fmt.Sprintf("%d", v1Asset.FileSize)
	}

	if len(properties) > 0 {
		asset.Properties = properties
	}

	return asset, nil
}

// MigrateTag converts a v1 Tag to v2 Tag
func (m *DataMigration) MigrateTag(v1Tag *v1entities.Tag) (*v2entities.Tag, error) {
	return &v2entities.Tag{
		Name: v1Tag.Name,
	}, nil
}

// MigrateAllProjects migrates all projects from v1 to v2
func (m *DataMigration) MigrateAllProjects(fs v2entities.LibFS) error {
	var v1Projects []*v1entities.Project

	// Load all v1 projects
	if err := m.v1DB.Preload("Assets").Preload("Tags").Find(&v1Projects).Error; err != nil {
		return fmt.Errorf("failed to load v1 projects: %w", err)
	}

	// Migrate each project
	for _, v1Project := range v1Projects {
		// Convert project to root asset
		rootAsset, err := m.MigrateProject(v1Project, fs)
		if err != nil {
			return fmt.Errorf("failed to migrate project %s: %w", v1Project.Name, err)
		}

		// Save root asset
		if err := m.v2DB.Create(rootAsset).Error; err != nil {
			return fmt.Errorf("failed to save migrated project %s: %w", v1Project.Name, err)
		}

		// Migrate project assets
		for _, v1Asset := range v1Project.Assets {
			nestedAsset, err := m.MigrateAsset(&v1Asset, fs, &rootAsset.ID)
			if err != nil {
				return fmt.Errorf("failed to migrate asset %s: %w", v1Asset.Name, err)
			}

			// Save nested asset
			if err := m.v2DB.Create(nestedAsset).Error; err != nil {
				return fmt.Errorf("failed to save migrated asset %s: %w", v1Asset.Name, err)
			}
		}
	}

	return nil
}

// MigrateAllTags migrates all tags from v1 to v2
func (m *DataMigration) MigrateAllTags() error {
	var v1Tags []*v1entities.Tag

	// Load all v1 tags
	if err := m.v1DB.Find(&v1Tags).Error; err != nil {
		return fmt.Errorf("failed to load v1 tags: %w", err)
	}

	// Migrate each tag
	for _, v1Tag := range v1Tags {
		v2Tag, err := m.MigrateTag(v1Tag)
		if err != nil {
			return fmt.Errorf("failed to migrate tag %s: %w", v1Tag.Name, err)
		}

		// Save tag (ignore duplicates)
		if err := m.v2DB.FirstOrCreate(v2Tag, v2entities.Tag{Name: v2Tag.Name}).Error; err != nil {
			return fmt.Errorf("failed to save migrated tag %s: %w", v1Tag.Name, err)
		}
	}

	return nil
}

// MigrateAll performs a complete migration from v1 to v2
func (m *DataMigration) MigrateAll(fs v2entities.LibFS) error {
	// Migrate tags first (they're referenced by projects)
	if err := m.MigrateAllTags(); err != nil {
		return fmt.Errorf("failed to migrate tags: %w", err)
	}

	// Migrate projects and assets
	if err := m.MigrateAllProjects(fs); err != nil {
		return fmt.Errorf("failed to migrate projects: %w", err)
	}

	return nil
}
