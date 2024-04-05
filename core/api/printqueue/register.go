package printqueue

import (
	"github.com/labstack/echo/v4"
)

func Register(e *echo.Group) {
	e.GET("/status", enqueue)
	e.POST("/enqueue", enqueue)
}
