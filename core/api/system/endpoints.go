package system

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"
	"sort"
	"strings"

	"github.com/eduardooliveira/stLib/core/data/database"
	"github.com/eduardooliveira/stLib/core/events"
	"github.com/eduardooliveira/stLib/core/processing"
	"github.com/eduardooliveira/stLib/core/runtime"
	"github.com/eduardooliveira/stLib/core/system"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/maps"
)

type void struct{}

func paths(c echo.Context) error {

	rtn := make(map[string]void, 0)
	filepath.WalkDir(runtime.Cfg.Library.Path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}

		path = strings.TrimLeft(path, runtime.Cfg.Library.Path)
		projectsPath := filepath.Clean(fmt.Sprintf("/%s", filepath.Dir(path)))
		projectName := filepath.Base(path)

		if p, err := database.GetProjectByPathAndName(projectsPath, projectName); err == nil && p.UUID != "" {
			return nil
		}

		rtn[path] = void{}

		return nil
	})
	s := maps.Keys(rtn)
	sort.Slice(s, func(i, j int) bool {
		return len(s[i]) < len(s[j])
	})
	return c.JSON(http.StatusOK, s)
}

func settings(c echo.Context) error {
	return c.JSON(http.StatusOK, runtime.Cfg)
}

func saveSettings(c echo.Context) error {
	cfg := &runtime.Config{}
	if err := c.Bind(cfg); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := runtime.SaveConfig(cfg); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, cfg)

}

func runDiscovery(c echo.Context) error {
	go processing.Run(runtime.Cfg.Library.Path)
	return c.NoContent(http.StatusOK)
}

func subscribe(c echo.Context) error {

	session := c.Param("session")
	if session == "" {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("no session provided").Error())
	}

	err := events.Subscribe(session, "system.state", system.GetEventPublisher())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func unSubscribe(c echo.Context) error {

	session := c.Param("session")
	if session == "" {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("no session provided").Error())
	}

	events.UnSubscribe(session, "system.state")

	return c.NoContent(http.StatusOK)
}
