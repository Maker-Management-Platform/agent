package processing

import (
	"github.com/eduardooliveira/stLib/core/entities"
)

type processableAsset struct {
	name    string
	label   string
	project *entities.Project
	asset   *entities.ProjectAsset
	parent  *entities.ProjectAsset
	origin  string
}

func (p *processableAsset) Name() string {
	return p.name
}
func (p *processableAsset) Project() *entities.Project {
	return p.project
}
func (p *processableAsset) Asset() *entities.ProjectAsset {
	return p.asset
}
