package repository

import (
	"context"

	"github.com/eduardooliveira/stLib/core/entities"
	"gorm.io/gorm"
)

// ProjectRepository defines the interface for project data access
type ProjectRepository interface {
	FindByID(ctx context.Context, id uint) (*entities.Project, error)
	FindAll(ctx context.Context) ([]*entities.Project, error)
	Create(ctx context.Context, project *entities.Project) error
	Update(ctx context.Context, project *entities.Project) error
	Delete(ctx context.Context, id uint) error
	FindByPath(ctx context.Context, path string) (*entities.Project, error)
}

// AssetRepository defines the interface for asset data access
type AssetRepository interface {
	FindByID(ctx context.Context, id uint) (*entities.Asset, error)
	FindAll(ctx context.Context) ([]*entities.Asset, error)
	FindByProjectID(ctx context.Context, projectID uint) ([]*entities.Asset, error)
	Create(ctx context.Context, asset *entities.Asset) error
	Update(ctx context.Context, asset *entities.Asset) error
	Delete(ctx context.Context, id uint) error
}

// TagRepository defines the interface for tag data access
type TagRepository interface {
	FindByID(ctx context.Context, id uint) (*entities.Tag, error)
	FindAll(ctx context.Context) ([]*entities.Tag, error)
	FindByName(ctx context.Context, name string) (*entities.Tag, error)
	Create(ctx context.Context, tag *entities.Tag) error
	Update(ctx context.Context, tag *entities.Tag) error
	Delete(ctx context.Context, id uint) error
}

// Repositories aggregates all repository interfaces
type Repositories struct {
	Projects ProjectRepository
	Assets   AssetRepository
	Tags     TagRepository
}

// NewRepositories creates a new repository collection
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Projects: NewProjectRepository(db),
		Assets:   NewAssetRepository(db),
		Tags:     NewTagRepository(db),
	}
}
