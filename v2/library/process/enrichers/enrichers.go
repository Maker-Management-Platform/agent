package enrichers

import (
	"context"

	"github.com/eduardooliveira/stLib/v2/library/entities"
)

type Enricher interface {
	Enrich(ctx context.Context, asset *entities.Asset) error
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
