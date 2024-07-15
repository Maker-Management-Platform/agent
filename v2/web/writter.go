package web

import (
	"github.com/a-h/templ"
	"github.com/eduardooliveira/stLib/v2/web/comp"
	"github.com/labstack/echo/v4"
)

type ResponseModel struct {
	Ctx          echo.Context
	S            int
	WrapperModel comp.WrapperModel
	IsFragment   bool
	PushState    string
}

func Render(rm ResponseModel) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	c := rm.WrapperModel.Main

	if rm.Ctx.Request().Header.Get("HX-Request") != "true" && !rm.IsFragment {
		c = comp.WrapperComponent(rm.WrapperModel)
	}

	if rm.PushState != "" {
		rm.Ctx.Response().Header().Set("HX-Push-Url", rm.PushState)
	}

	if err := c.Render(rm.Ctx.Request().Context(), buf); err != nil {
		return err
	}

	return rm.Ctx.HTML(rm.S, buf.String())
}

func Error(ctx echo.Context, statusCode int, message string) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := comp.WrapperComponent(comp.WrapperModel{Main: comp.Error(message)}).
		Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}
