package comp

import (
	"github.com/a-h/templ"
	"github.com/eduardooliveira/stLib/v2/config"
	"github.com/eduardooliveira/stLib/v2/library/entities"
)

type IndexModel struct {
	Asset       *entities.Asset
	AssetTypes  []config.AssetType
	KindFilter  KindFilterModel
	Main        templ.Component
	Pagination  PaginationModel
	SearchModel SearchModel
}

type ListModel struct {
	Asset      *entities.Asset
	Pagination PaginationModel
}

type AssetCardModel struct {
	Asset *entities.Asset
}

type DetailsModel struct {
	Asset *entities.Asset
}

type EditModel struct {
	Asset  *entities.Asset
	Errors map[string]string
	Action string
}

type NewModel struct {
	ParentID  string
	TempFiles []string
	Errors    map[string]string
}

type PaginationModel struct {
	CurrentPage int
	TotalPages  int
	OOB         bool
	Fields      map[string]string
}

type KindFilterModel struct {
	AssetTypes []config.AssetType
	Selected   string
}

type SearchModel struct {
	Search bool   `in:"query=search"`
	Global bool   `in:"query=global"`
	Parent string `in:"query=parent"`
	Name   string `in:"query=name"`
	Tags   string `in:"query=tags"`
}
