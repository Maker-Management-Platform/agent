package web

import (
	"net/http"

	"github.com/eduardooliveira/stLib/v2/library/web/comp"
	"github.com/eduardooliveira/stLib/v2/web"
)

func (h webHandler) sidebarHandler(r *http.Request) web.ResponseModel {
	return web.ResponseModel{
		S:          http.StatusOK,
		Component:  comp.View3D(),
		IsFragment: true,
	}
}
