package enrichers

import (
	"slices"

	"github.com/eduardooliveira/stLib/v2/library/entities"
)

type Enricher interface {
	Enrich(asset *entities.Asset) func() error
}

var extensions = []string{
	".gcode",
}

var enrichers = map[string]Enricher{}

func IsEnrichable(asset *entities.Asset) bool {
	return slices.Contains(extensions, *asset.Extension)
}

func GetEnricher(asset *entities.Asset) Enricher {
	return enrichers[*asset.Extension]
}

func Get(asset *entities.Asset) (Enricher, bool) {
	enricher, ok := enrichers[*asset.Extension]
	return enricher, ok
}
