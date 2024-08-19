package renderers

import (
	"github.com/eduardooliveira/stLib/v2/library/entities"
)

type Renderer interface {
	Render(asset *entities.Asset) (*entities.Asset, error)
}

var renderers = map[string]Renderer{}

func Init() error {
	renderers = map[string]Renderer{
		".gcode": &gCodeRenderer{},
		".stl":   NewSTLRenderer(),
	}
	return nil
}

func Get(asset *entities.Asset) (Renderer, bool) {
	renderer, ok := renderers[*asset.Extension]
	return renderer, ok
}
