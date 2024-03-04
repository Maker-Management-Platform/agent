package processing

import (
	"github.com/eduardooliveira/stLib/core/entities"
)

type processableAsset struct {
	name      string
	path      string
	project   *entities.Project
	asset     *entities.ProjectAsset
	parent    *entities.ProjectAsset
	generated bool
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
