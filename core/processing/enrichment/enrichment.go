package enrichment

import "github.com/eduardooliveira/stLib/core/entities"

type Enrichable interface {
	Asset() *entities.ProjectAsset
	Project() *entities.Project
}

type Renderer interface {
	Render(Enrichable) (string, error)
}

type Parser interface {
	Parse(Enrichable) error
}

type Extractor interface {
	Extract(Enrichable) ([]string, error)
}
