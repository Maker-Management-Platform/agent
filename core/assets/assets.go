package assets

import (
	"github.com/labstack/echo/v4"
)

func Register(e *echo.Group) {

	e.POST("/:id/delete", deleteAsset)
	//e.POST("/:id", save) not in use
	e.POST("", new)
}
