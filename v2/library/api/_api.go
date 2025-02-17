package api

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"path"

	"github.com/eduardooliveira/stLib/v2/config"
	"github.com/eduardooliveira/stLib/v2/library/repo"
	"github.com/eduardooliveira/stLib/v2/utils"
	"github.com/labstack/echo/v4"
)

type API struct {
	r *repo.AssetRepo
	l *slog.Logger
}

func New(e echo.Group, r *repo.AssetRepo) error {
	api := &API{
		r: r,
		l: slog.With("module", "api"),
	}
	if err := utils.CreateFolder(path.Join(config.Cfg.Core.DataFolder, "temp")); err != nil {
		return err
	}
	e.GET("/api/version", api.version)      // octoprint / prusa connect
	e.POST("/api/files/local", api.upload)  // octoprint / prusa connect
	e.GET("/server/info", api.info)         // klipper
	e.POST("/api/files/upload", api.upload) // klipper
	return nil
}

func (a API) version(c echo.Context) error {
	return c.JSON(http.StatusOK, struct {
		API    string `json:"api"`
		Server string `json:"server"`
		Text   string `json:"text"`
	}{
		API:    "xxx",
		Server: "xxx",
		Text:   "OctoPrint xx",
	})
}

func (a API) info(c echo.Context) error {
	return c.JSON(http.StatusOK, struct {
		State           string `json:"state"`
		StateMessage    string `json:"state_message"`
		Hostname        string `json:"hostname"`
		SoftwareVersion string `json:"software_version"`
		CPUInfo         string `json:"cpu_info"`
		KlipperPath     string `json:"klipper_path"`
		PythonPath      string `json:"python_path"`
		LogFile         string `json:"log_file"`
		ConfigFile      string `json:"config_file"`
	}{
		State:           "ready",
		StateMessage:    "Printer is ready",
		Hostname:        "mmp",
		SoftwareVersion: "v0.9.xxx",
		CPUInfo:         "xxx",
		KlipperPath:     "/root/klipper", // This are mock values, do not run stuff as root....
		PythonPath:      "/root/klippy-env/bin/python",
		LogFile:         "/tmp/klippy.log",
		ConfigFile:      "/root/printer.cfg",
	})
}

func (a API) upload(c echo.Context) error {

	form, err := c.MultipartForm()
	if err != nil {
		a.l.Error("upload", "err", err)
		return c.NoContent(http.StatusBadRequest)
	}

	files := form.File["file"]

	if len(files) == 0 {
		a.l.Error("upload", "err", "no files")
		return c.NoContent(http.StatusBadRequest)
	}
	name := files[0].Filename

	// Source
	src, err := files[0].Open()
	if err != nil {
		a.l.Error("upload", "err", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(path.Join(config.Cfg.Core.DataFolder, "temp", name))
	if err != nil {
		a.l.Error("upload", "err", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		a.l.Error("upload", "err", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
