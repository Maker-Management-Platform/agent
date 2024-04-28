package api

import "github.com/labstack/echo/v4"

func Register(e *echo.Group) {

	e.POST("", new)
	e.GET("", index)
	e.GET("/:uuid", show)
	e.GET("/:uuid/stream", stream)
	e.POST("/:uuid", edit)
	e.POST("/:uuid/delete", deleteHandler)
	e.GET("/:uuid/send/:id", sendHandler)
	e.GET("/:uuid/subscribe/:session", subscribe)
	e.GET("/:uuid/unsubscribe/:session", unsubscribe)
	e.POST("/test", testConnection)
}
