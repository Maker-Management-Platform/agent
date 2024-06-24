package renderers

import (
	"slices"

	"github.com/eduardooliveira/stLib/v2/library/entities"
)

type Renderer interface {
	Render(asset entities.Asset, cb OnRenderCallback) func() error
}
type OnRenderCallback func(*entities.Asset, string, string) error

var extensions = []string{
	".gcode",
	".stl",
}

var renderers = map[string]Renderer{}

func Init() error {
	renderers = map[string]Renderer{
		".gcode": &gCodeRenderer{},
		".stl":   NewSTLRenderer(),
	}
	return nil
}

func IsRenderable(asset *entities.Asset) bool {
	return slices.Contains(extensions, *asset.Extension)
}

func GetRenderer(asset *entities.Asset) Renderer {
	return renderers[*asset.Extension]
}

func Get(asset *entities.Asset) (Renderer, bool) {
	renderer, ok := renderers[*asset.Extension]
	return renderer, ok
}
