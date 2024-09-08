package web

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/eduardooliveira/stLib/v2/web/comp"
)

func (h webHandler) indexHandler(r *http.Request) ResponseModel {
	return ResponseModel{
		Status: http.StatusOK,
		Component: comp.WrapperComponent(comp.WrapperModel{
			Main: templ.NopComponent,
		}),
		IsFragment: true,
	}
}
