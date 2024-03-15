package events

import (
	"github.com/labstack/echo/v4"
)

func Register(e *echo.Group) {

	group := e
	group.GET("", index)
}
