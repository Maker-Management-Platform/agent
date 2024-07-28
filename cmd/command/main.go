package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/eduardooliveira/stLib/v2/config"
	"github.com/eduardooliveira/stLib/v2/database"
	"github.com/eduardooliveira/stLib/v2/library"
	"github.com/eduardooliveira/stLib/v2/utils"
	"github.com/eduardooliveira/stLib/v2/web/mw"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
	var dataFolder string
	flag.StringVar(&dataFolder, "data-folder", "", "Data folder")
	flag.BoolVar(&utils.IsDocker, "docker", false, "Running in docker")
	flag.Parse()

	if dataFolder == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Error getting user home directory: %v", err)
		}
		dataFolder = filepath.Join(home, ".mmp")
	}

	if _, err := os.Stat(dataFolder); os.IsNotExist(err) {
		if err := os.Mkdir(dataFolder, 0755); err != nil {
			log.Fatalf("Error creating data folder: %v", err)
		}
	}

	if err := utils.CreateFolder(filepath.Join(dataFolder, "img")); err != nil {
		log.Fatalf("Error creating data folder: %v", err)
	}

	if err := config.Init(dataFolder); err != nil {
		log.Fatalf("Error initializing config: %v", err)
	}

	if err := database.Init(); err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	server := echo.New()
	server.Use(mw.CtxMiddleware())
	server.Use(middleware.CORS())
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())

	l, err := library.New(server.Group(""))
	if err != nil {
		log.Fatalf("Error initializing library: %v", err)
	}

	l.ScanAsync()

	slog.Info("Starting agent")
	log.Fatal(server.Start(fmt.Sprintf(":%d", config.Cfg.Server.Port)))
}
