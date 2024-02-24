package entities

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/eduardooliveira/stLib/core/utils"
	"github.com/google/uuid"
)

type Project struct {
	UUID           string                   `json:"uuid" toml:"uuid" form:"uuid" query:"uuid" gorm:"primaryKey"`
	Name           string                   `json:"name" toml:"name" form:"name" query:"name"`
	Description    string                   `json:"description,omitempty" toml:"description" form:"description" query:"description"`
	Path           string                   `json:"path,omitempty" toml:"path" form:"path" query:"path"`
	ExternalLink   string                   `json:"external_link,omitempty" toml:"external_link" form:"external_link" query:"external_link"`
	Assets         map[string]*ProjectAsset `json:"-" toml:"-" form:"assets" query:"assets" gorm:"-"`
	Tags           []*Tag                   `json:"tags,omitempty" toml:"tags" form:"tags" query:"tags" gorm:"many2many:project_tags"`
	DefaultImageID string                   `json:"default_image_id,omitempty" toml:"default_image_id" form:"default_image_id" query:"default_image_id"`
	Initialized    bool                     `json:"initialized,omitempty" toml:"initialized" form:"initialized" query:"initialized"`
}

func (p *Project) FullPath() string {
	return filepath.Clean(path.Join(p.Path, p.Name))
}

func NewProjectFromPath(path string) *Project {
	if p := tryLoadFromFile(path); p != nil {
		return p
	}

	project := NewProject()
	project.Path = filepath.Clean(fmt.Sprintf("/%s", filepath.Dir(path)))
	project.Name = filepath.Base(path)
	return project
}

func tryLoadFromFile(path string) *Project {
	p := NewProject()
	_, err := toml.DecodeFile(utils.ToLibPath(fmt.Sprintf("%s/.project.stlib", path)), &p)
	if err != nil {
		return nil
	}
	p.Initialized = true
	return p
}

func NewProject() *Project {
	project := &Project{
		UUID:        uuid.New().String(),
		Initialized: false,
		Tags:        make([]*Tag, 0),
		Assets:      make(map[string]*ProjectAsset),
	}
	return project
}
