package system

import (
	"github.com/labstack/echo/v4"
)

var group *echo.Group

func Register(e *echo.Group) {
	group = e
	group.GET("/paths", paths)
	group.GET("/settings", settings)
	group.POST("/settings", saveSettings)
	group.GET("/discovery", runDiscovery)
	group.GET("/events/subscribe/:session", subscribe)
	group.GET("/events/unsubscribe/:session", unSubscribe)
}
