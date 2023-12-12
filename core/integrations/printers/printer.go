package printers

import "github.com/labstack/echo/v4"

var group *echo.Group

func Register(e *echo.Group) {

	group = e
	//group.GET("", index)
	group.POST("", save)
}
