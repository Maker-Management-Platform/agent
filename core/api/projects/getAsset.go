package projects

import (
	"errors"
	"log"
	"net/http"
	"path"

	"github.com/eduardooliveira/stLib/core/data/database"
	"github.com/eduardooliveira/stLib/core/utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func getAsset(c echo.Context) error {
	project, err := database.GetProject(c.Param("uuid"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	asset, err := database.GetProjectAsset(project.UUID, c.Param("id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var assetPath string
	if asset.Origin == "fs" {
		assetPath = utils.ToLibPath(path.Join(project.FullPath(), asset.Name))
	} else {
		assetPath = utils.ToAssetsPath(project.UUID, asset.Name)
	}

	if c.QueryParam("download") != "" {
		return c.Attachment(assetPath, asset.Name)

	}

	return c.Inline(assetPath, asset.Name)
}
