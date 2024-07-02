package mw

import (
	"github.com/eduardooliveira/stLib/v2/web/helpers"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slog"
)

type contextMiddleware struct {
	echo.Context
	uh     *helpers.URLHelper
	isHtmx bool
}

func CtxMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &contextMiddleware{
				Context: c,
				uh:      helpers.NewURLHelper(c.Request()),
				isHtmx:  c.Request().Header.Get("HX-Request") == "true",
			}
			return next(cc)
		}
	}
}

func URLHelper(ctx echo.Context) *helpers.URLHelper {
	uc, ok := ctx.(*contextMiddleware)
	if ok {
		return uc.uh
	}
	slog.Error("not a contextMiddleware")
	return nil
}
