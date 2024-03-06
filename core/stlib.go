package stlib

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/eduardooliveira/stLib/core/api/projects"
	"github.com/eduardooliveira/stLib/core/api/system"
	"github.com/eduardooliveira/stLib/core/api/tags"
	"github.com/eduardooliveira/stLib/core/api/tempfiles"
	"github.com/eduardooliveira/stLib/core/data/database"
	"github.com/eduardooliveira/stLib/core/downloader"
	"github.com/eduardooliveira/stLib/core/events"
	"github.com/eduardooliveira/stLib/core/integrations/printers"
	"github.com/eduardooliveira/stLib/core/integrations/slicer"
	"github.com/eduardooliveira/stLib/core/processing"
	"github.com/eduardooliveira/stLib/core/runtime"
	"github.com/eduardooliveira/stLib/core/state"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Run() {
	if runtime.Cfg.Core.Log.EnableFile && runtime.Cfg.Core.Log.Path != "" {
		f, err := os.OpenFile(filepath.Join(runtime.Cfg.Core.Log.Path, "log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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

	err = state.LoadAssetTypes()
	if err != nil {
		log.Fatal("error loading assetTypes", err)
	}
	go processing.Run(runtime.Cfg.Library.Path)
	go processing.RunTempDiscovery()
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
	events.Register(api.Group("/events"))
	projects.Register(api.Group("/projects"))
	tags.Register(api.Group("/tags"))
	tempfiles.Register(api.Group("/tempfiles"))
	printers.Register(api.Group("/printers"))
	downloader.Register(api.Group("/downloader"))
	system.Register(api.Group("/system"))
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", runtime.Cfg.Server.Port)))
}
