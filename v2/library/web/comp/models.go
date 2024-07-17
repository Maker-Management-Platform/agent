package comp

import (
	"github.com/a-h/templ"
	"github.com/eduardooliveira/stLib/v2/config"
	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/web/helpers"

	corecomp "github.com/eduardooliveira/stLib/v2/web/comp"
)

type IndexModel struct {
	Asset      *entities.Asset
	AssetTypes []config.AssetType
	Main       templ.Component
}

type ListModel struct {
	Asset      *entities.Asset
	Pagination corecomp.PaginationModel
}

type AssetCardModel struct {
	Asset *entities.Asset
	UH    *helpers.URLHelper
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
