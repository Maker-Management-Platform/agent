package processing

import (
	"github.com/eduardooliveira/stLib/core/entities"
)

type ProcessableAsset struct {
	Name    string
	Label   string
	Project *entities.Project
	Asset   *entities.ProjectAsset
	Parent  *entities.ProjectAsset
	Origin  string
}

func (p *ProcessableAsset) GetName() string {
	return p.Name
}
func (p *ProcessableAsset) GetProject() *entities.Project {
	return p.Project
}
func (p *ProcessableAsset) GetAsset() *entities.ProjectAsset {
	return p.Asset
}
