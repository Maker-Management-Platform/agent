package projects

import (
	"github.com/eduardooliveira/stLib/core/api/projects/assets"
	"github.com/labstack/echo/v4"
)

var group *echo.Group

func Register(e *echo.Group) {

	group = e
	group.GET("", index)
	group.GET("/list", list)
	group.GET("/:uuid", show)
	group.GET("/:uuid/discover", discoverHandler)
	group.POST("/:uuid", save)
	group.POST("/:uuid/move", moveHandler)
	group.POST("/:uuid/image", setMainImageHandler)
	group.POST("/:uuid/delete", deleteHandler)
	group.POST("", new)

	group.GET("/:uuid/assets", assets.List)
	group.POST("/:uuid/assets", assets.New)
	group.GET("/:uuid/assets/:id", assets.Get)
	group.POST("/:uuid/assets/:id/delete", assets.Delete)
}
