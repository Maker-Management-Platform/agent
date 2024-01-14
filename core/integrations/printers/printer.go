package printers

import (
	"log"
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
	group.POST("/:uuid", edit)
	group.POST("/:uuid/delete", deleteHandler)
	group.GET("/:uuid/send/:id", sendHandler)
	group.POST("/test", testConnection)

	go checkConnection()
}

func checkConnection() {
	for {
		log.Println("Checking printer connectivity")
		for _, p := range state.Printers {
			switch p.Type {
			case "klipper":
				klipper.ConntectionStatus(p)
			}
		}
		time.Sleep(10 * time.Second)
	}

}
