package printqueue

import (
	"github.com/labstack/echo/v4"
)

func Register(e *echo.Group) {
	e.GET("/jobs", getPrintJobs)
	e.GET("/jobs/:uuid/cancel", cancel)
	e.GET("/move", move)
	e.POST("/enqueue", enqueue)

	e.GET("/subscribe/:session", subscribe)
	e.GET("/unsubscribe/:session", unSubscribe)
}
