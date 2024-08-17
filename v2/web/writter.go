package web

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/eduardooliveira/stLib/v2/web/comp"
	"github.com/labstack/echo/v4"
)

type ResponseModel struct {
	Ctx          echo.Context
	S            int
	Error        error
	Component    templ.Component
	OOB          []templ.Component
	WrapperModel comp.WrapperModel
	IsFragment   bool
	PushState    string
	Events       []string
}

type RenderRequest func(r *http.Request) ResponseModel

func R(rr RenderRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rm := rr(r)

		isFragment := r.Header.Get("Hx-Request") == "true" || rm.IsFragment

		if rm.Error != nil {
			if isFragment {
				comp.Error(rm.Error.Error()).Render(r.Context(), w)
			} else {
				comp.WrapperComponent(comp.WrapperModel{Main: comp.Error(rm.Error.Error())}).Render(r.Context(), w)
			}
			return
		}
		slog.Info("Rendering", slog.Any("responseModel", rm))
		if rm.PushState != "" {
			w.Header().Add("HX-Push-Url", rm.PushState)
		}

		if rm.Events != nil && len(rm.Events) > 0 {
			w.Header().Add("HX-Trigger", strings.Join(rm.Events, ", "))
		}

		w.WriteHeader(rm.S)
		fComponent := rm.Component
		if !isFragment {
			fComponent = comp.WrapperComponent(comp.WrapperModel{Main: rm.Component})
		}
		if err := fComponent.Render(r.Context(), w); err != nil {
			if isFragment {
				comp.Error(err.Error()).Render(r.Context(), w)
			} else {
				comp.WrapperComponent(comp.WrapperModel{Main: comp.Error(err.Error())}).Render(r.Context(), w)
			}
		}
		if isFragment && rm.OOB != nil {
			for _, c := range rm.OOB {
				c.Render(r.Context(), w)
			}
		}

	}
}
