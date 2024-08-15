package comp

import (
	"github.com/a-h/templ"
	"github.com/eduardooliveira/stLib/v2/web/helpers"
)

type WrapperModel struct {
	Main   templ.Component
	AsideR templ.Component
	Other  map[string]templ.Component
}

type PaginationModel struct {
	UH          *helpers.URLHelper
	CurrentPage int
	TotalPages  int
	Target      string
}
