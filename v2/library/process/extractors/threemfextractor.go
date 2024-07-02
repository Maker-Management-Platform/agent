package extractors

import "github.com/eduardooliveira/stLib/v2/library/entities"

type ThreeMFExtractor struct {
}

func (t *ThreeMFExtractor) Extract(asset *entities.Asset) ([]*entities.Asset, error) {
	return []*entities.Asset{asset}, nil
}
