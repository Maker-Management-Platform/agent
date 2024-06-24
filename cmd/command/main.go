package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/eduardooliveira/stLib/v2/config"
	"github.com/eduardooliveira/stLib/v2/database"
	"github.com/eduardooliveira/stLib/v2/library"
	"github.com/eduardooliveira/stLib/v2/library/discovery"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/sync/errgroup"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
	var dataFolder string
	flag.StringVar(&dataFolder, "data-folder", "", "Data folder")
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
	if _, err := os.Stat(filepath.Join(dataFolder, "img")); os.IsNotExist(err) {
		if err := os.Mkdir(filepath.Join(dataFolder, "img"), 0755); err != nil {
			log.Fatalf("Error creating data folder: %v", err)
		}
	}

	if err := config.Init(dataFolder); err != nil {
		log.Fatalf("Error initializing config: %v", err)
	}

	if err := database.Init(); err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	server := echo.New()

	server.Use(middleware.CORS())
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())

	if err := library.Init(*server.Group("/lib")); err != nil {
		log.Fatalf("Error initializing library: %v", err)
	}

	eg := errgroup.Group{} //maybe this should be initialized in the library.Init() function
	if len(config.Cfg.Library.Paths) == 0 {
		slog.Warn("No library paths configured")
	}
	for _, path := range config.Cfg.Library.Paths {
		eg.Go(discovery.New(path).Run)
	}

	eg.Go(func() error {
		return server.Start(fmt.Sprintf(":%d", config.Cfg.Server.Port))
	})

	eg.Go(func() error {
		time.Sleep(time.Minute)

		return nil
	})

	slog.Info("Starting agent")
	if err := eg.Wait(); err != nil {
		log.Fatalf("Error running agent: %v", err)
	}
}
