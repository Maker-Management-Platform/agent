package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/eduardooliveira/stLib/frontend"
	"github.com/eduardooliveira/stLib/v2/config"
	"github.com/eduardooliveira/stLib/v2/database"
	"github.com/eduardooliveira/stLib/v2/library"
	"github.com/eduardooliveira/stLib/v2/utils"
	"github.com/eduardooliveira/stLib/v2/web"
	"github.com/go-chi/chi/v5"
	cmw "github.com/go-chi/chi/v5/middleware"
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

	if err := config.Load(dataFolder); err != nil {
		log.Fatalf("Error initializing config: %v", err)
	}

	if err := database.Init(); err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	r := chi.NewRouter()
	r.Use(cmw.Logger)
	r.Use(cmw.Recoverer)

	r.Handle("/dist/*", http.FileServerFS(frontend.FS))

	webH, err := web.New()
	if err != nil {
		log.Fatalf("Error initializing web: %v", err)
	}
	r.Mount("/", webH)

	l, libH, _, err := library.New()
	if err != nil {
		log.Fatalf("Error initializing library: %v", err)
	}

	r.Mount("/lib", libH)

	l.ScanAsync(context.Background())

	slog.Info("Starting agent")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Cfg.Server.Port), r))
}
