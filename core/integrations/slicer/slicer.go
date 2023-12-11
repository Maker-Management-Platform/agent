package slicer

import (
	"github.com/labstack/echo/v4"
)

var group *echo.Group

func Register(e *echo.Group) {
	group = e
	group.GET("/api/version", version)
	group.POST("/api/files/local", upload)
}
