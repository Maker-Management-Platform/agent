package main

import (
	"context"
	"flag"
	"fmt"
	"io/fs"
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
	"github.com/go-chi/chi/v5"
	cmw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	assets, err := fs.Sub(frontend.FS, "dist")
	if err != nil {
		log.Fatalf("Error getting asset: %v", err)
	}
	r.Handle("/assets/*", http.FileServerFS(assets))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, frontend.FS, "/dist/index.html")
	})

	l, libH, _, err := library.New()
	if err != nil {
		log.Fatalf("Error initializing library: %v", err)
	}
	r.Route("/api", func(r chi.Router) {
		r.Mount("/lib", libH)
	})

	l.ScanAsync(context.Background())

	slog.Info("Starting agent")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Cfg.Server.Port), r))
}
