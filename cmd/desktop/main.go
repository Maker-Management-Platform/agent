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
	"github.com/go-chi/chi/v5"
	cmw "github.com/go-chi/chi/v5/middleware"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
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

	if err := config.Load(dataFolder); err != nil {
		log.Fatalf("Error initializing config: %v", err)
	}

	if err := database.Init(); err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	//eg := errgroup.Group{}

	r := chi.NewRouter()
	r.Use(cmw.Logger)
	r.Use(cmw.Recoverer)

	l, libH, _, err := library.New()
	if err != nil {
		log.Fatalf("Error initializing library: %v", err)
	}

	r.Mount("/lib", libH)

	l.ScanAsync(context.Background())

	slog.Info("Starting agent")
	go func() {

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Cfg.Server.Port), r))
	}()

	wails.Run(&options.App{
		Title:  "mmp",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: frontend.FS,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		// Windows platform specific options
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
			// DisableFramelessWindowDecorations: false,
			WebviewUserDataPath: "",
			ZoomFactor:          1.0,
		},
		// Mac platform specific options
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: false,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            false,
				UseToolbar:                 false,
				HideToolbarSeparator:       true,
			},
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   "MMP",
				Message: "",
			},
		},
	})

}
