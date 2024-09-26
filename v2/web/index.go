package web

import (
	"net/http"

	"github.com/eduardooliveira/stLib/v2/web/comp"
)

func (h webHandler) indexHandler(r *http.Request) ResponseModel {
	return ResponseModel{
		Status: http.StatusOK,
		Component: comp.WrapperComponent(comp.WrapperModel{
			Main: comp.Home(),
		}),
		IsFragment: true,
	}
}
