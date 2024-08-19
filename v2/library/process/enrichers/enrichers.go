package enrichers

import (
	"github.com/eduardooliveira/stLib/v2/library/entities"
)

type Enricher interface {
	Enrich(asset *entities.Asset) error
}

var enrichers = map[string]Enricher{}

func Init() error {
	enrichers = map[string]Enricher{
		".gcode": &gCodeEnricher{},
	}
	return nil
}

func Get(asset *entities.Asset) (Enricher, bool) {
	enricher, ok := enrichers[*asset.Extension]
	return enricher, ok
}
