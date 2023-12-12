package slicer

import (
	"github.com/labstack/echo/v4"
)

var group *echo.Group

func Register(e *echo.Group) {
	group = e
	group.GET("/api/version", version)      // octoprint / prusa connect
	group.POST("/api/files/local", upload)  // octoprint / prusa connect
	group.GET("/server/info", info)         // klipper
	group.POST("/api/files/upload", upload) // klipper
}
