package library

import (
	"github.com/eduardooliveira/stLib/v2/library/renderers"
	"github.com/eduardooliveira/stLib/v2/library/svc"
	"github.com/eduardooliveira/stLib/v2/library/web"
	"github.com/labstack/echo/v4"
)

func Init(e echo.Group) error {
	if err := renderers.Init(); err != nil {
		return err
	}

	if err := svc.Init(); err != nil {
		return err
	}

	if err := web.Init(e); err != nil {
		return err
	}

	return nil
}
