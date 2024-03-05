package tags

import (
	"github.com/labstack/echo/v4"
)

func Register(e *echo.Group) {
	e.GET("", index)
}
