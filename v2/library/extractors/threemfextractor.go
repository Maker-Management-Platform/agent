package extractors

import "github.com/eduardooliveira/stLib/v2/library/entities"

type ThreeMFExtractor struct {
}

func (t *ThreeMFExtractor) Extract(asset *entities.Asset, cb func([]*entities.Asset) error) func() error {
	return func() error {
		return cb([]*entities.Asset{asset})
	}
}
