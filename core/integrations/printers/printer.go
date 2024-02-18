package printers

import (
	"time"

	"github.com/eduardooliveira/stLib/core/integrations/klipper"
	"github.com/eduardooliveira/stLib/core/state"
	"github.com/labstack/echo/v4"
)

var group *echo.Group

func Register(e *echo.Group) {

	group = e
	group.POST("", new)
	group.GET("", index)
	group.GET("/:uuid", show)
	group.GET("/:uuid/stream", stream)
	group.GET("/:uuid/status", statusHandler)
	group.POST("/:uuid", edit)
	group.POST("/:uuid/delete", deleteHandler)
	group.GET("/:uuid/send/:id", sendHandler)
	group.GET("/:uuid/subscribe/:session", subscribe)
	group.GET("/:uuid/unsubscribe/:session", unSubscribe)
	group.POST("/test", testConnection)

	go checkConnection()
}

func checkConnection() {
	for {
		for _, p := range state.Printers {
			switch p.Type {
			case "klipper":
				klipper.ConnectionStatus(p)
			}
		}
		time.Sleep(10 * time.Second)
	}

}
