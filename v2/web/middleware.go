package web

import (
	"github.com/a-h/templ"
	"github.com/eduardooliveira/stLib/v2/web/components"
	"github.com/labstack/echo/v4"
)

func Render(ctx echo.Context, statusCode int, t components.Wrapper, htmx bool) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	c := t.MainContent

	if !htmx || ctx.Request().Header.Get("HX-Request") != "true" {
		c = components.WrapperComponent(t)
	} else {
		ctx.Response().Header().Set("HX-Push-Url", ctx.Request().URL.String())
	}

	if err := c.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

func Error(ctx echo.Context, statusCode int, message string) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := components.WrapperComponent(components.Wrapper{
		MainContent: components.Error(message),
	}).Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}
