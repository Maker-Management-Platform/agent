package stlib

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/eduardooliveira/stLib/core/assets"
	"github.com/eduardooliveira/stLib/core/data/database"
	"github.com/eduardooliveira/stLib/core/discovery"
	"github.com/eduardooliveira/stLib/core/downloader"
	"github.com/eduardooliveira/stLib/core/integrations/printers"
	"github.com/eduardooliveira/stLib/core/integrations/slicer"
	"github.com/eduardooliveira/stLib/core/projects"
	"github.com/eduardooliveira/stLib/core/runtime"
	"github.com/eduardooliveira/stLib/core/state"
	"github.com/eduardooliveira/stLib/core/system"
	"github.com/eduardooliveira/stLib/core/tags"
	"github.com/eduardooliveira/stLib/core/tempfiles"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Run() {

	if logPath := runtime.Cfg.LogPath; logPath != "" {
		f, err := os.OpenFile("stlib.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()
		wrt := io.MultiWriter(os.Stdout, f)
		log.SetOutput(wrt)
	}

	err := database.InitDatabase()
	if err != nil {
		log.Fatal("error initing database", err)
	}
	go discovery.Run(runtime.Cfg.LibraryPath)
	go discovery.RunTempDiscovery()
	err = state.LoadPrinters()
	if err != nil {
		log.Fatal("error loading printers", err)
	}

	fmt.Println("starting server...")
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	/*e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "frontend/dist",
		Index:  "index.html",
		Browse: false,
		HTML5:  true,
	}))*/

	slicer.Register(e.Group(""))

	api := e.Group("/api")

	projects.Register(api.Group("/projects"))
	tags.Register(api.Group("/tags"))
	assets.Register(api.Group("/assets"))
	tempfiles.Register(api.Group("/tempfiles"))
	printers.Register(api.Group("/printers"))
	downloader.Register(api.Group("/downloader"))
	system.Register(api.Group("/system"))
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", runtime.Cfg.Port)))
}
