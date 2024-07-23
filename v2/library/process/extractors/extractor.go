package extractors

import (
	"slices"

	"github.com/eduardooliveira/stLib/v2/library/entities"
)

type Extractor interface {
	Extract(asset *entities.Asset) ([]*entities.Asset, error)
	ExtractBundled(asset *entities.Asset) error
}

var extensions = []string{
	".3mf",
	".zip",
	".rar",
	".7z",
	".tar",
}

var extractors = map[string]Extractor{
	".3mf": &ThreeMFExtractor{},
	".zip": &StdExtractor{},
	".rar": &StdExtractor{},
	".7z":  &StdExtractor{},
	".tar": &StdExtractor{},
}

func IsExtractable(asset *entities.Asset) bool {
	return slices.Contains(extensions, *asset.Extension)
}

func GetExtractor(asset *entities.Asset) Extractor {
	return extractors[*asset.Extension]
}

func Get(asset *entities.Asset) (Extractor, bool) {
	e, ok := extractors[*asset.Extension]
	return e, ok
}
