package projects

import (
	"github.com/labstack/echo/v4"
)

var group *echo.Group

func Register(e *echo.Group) {

	group = e
	group.GET("", index)
	group.GET("/list", list)
	group.GET("/:uuid", show)
	group.GET("/:uuid/assets", showAssets)
	group.GET("/:uuid/assets/:id", getAsset)
	group.POST("/:uuid", save)
	group.POST("/:uuid/move", moveHandler)
	group.POST("/:uuid/image", setMainImageHandler)
	group.POST("/:uuid/delete", deleteHandler)
	group.POST("", new)
}
