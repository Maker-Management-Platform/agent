package enrichment

import (
	"context"
	"log"

	"github.com/eduardooliveira/stLib/core/entities"
	"github.com/eduardooliveira/stLib/core/processing/types"
)

type Extracted struct {
	Label string
	File  string
}

type Enrichable interface {
	GetAsset() *entities.ProjectAsset
	GetProject() *entities.Project
}

type Renderer interface {
	Render(types.ProcessableAsset) (string, error)
}

type Parser interface {
	Parse(types.ProcessableAsset) error
}

type Extractor interface {
	Extract(types.ProcessableAsset) ([]*Extracted, error)
}

var renderers = make(map[string][]Renderer, 0)
var parsers = make(map[string][]Parser, 0)
var extractors = make(map[string][]Extractor, 0)

func init() {
	renderers[".stl"] = []Renderer{NewSTLRenderer()}
	renderers[".gcode"] = []Renderer{NewGCodeRenderer()}
	parsers[".gcode"] = []Parser{NewGCodeParser()}
	extractors[".3mf"] = []Extractor{New3MFExtractor()}
}

func EnrichAsset(ctx context.Context, p *types.ProcessableAsset) ([]*types.ProcessableAsset, error) {
	rtn := make([]*types.ProcessableAsset, 0)
	if rs, ok := renderers[p.Asset.Extension]; ok {
		for _, r := range rs {
			file, err := r.Render(*p)
			if err != nil {
				return nil, err
			}
			rtn = append(rtn, &types.ProcessableAsset{
				Name:    file,
				Project: p.Project,
				Origin:  "render",
			})
		}
	}
	if es, ok := extractors[p.Asset.Extension]; ok {
		for _, e := range es {
			extracted, err := e.Extract(*p)
			if err != nil {
				return nil, err
			}
			for _, e := range extracted {
				rtn = append(rtn, &types.ProcessableAsset{
					Name:    e.File,
					Label:   e.Label,
					Project: p.Project,
					Origin:  "extract",
				})
			}
		}
	}
	if ps, ok := parsers[p.Asset.Extension]; ok {
		for _, parser := range ps {
			if err := parser.Parse(*p); err != nil {
				log.Println(err)
			}
		}
	}
	return rtn, nil
}
