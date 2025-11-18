package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eduardooliveira/stLib/core/integration"
	v1stlib "github.com/eduardooliveira/stLib/core"
	v2config "github.com/eduardooliveira/stLib/v2/config"
	v2database "github.com/eduardooliveira/stLib/v2/database"
	v2library "github.com/eduardooliveira/stLib/v2/library"
	"github.com/go-chi/chi/v5"
	cmw "github.com/go-chi/chi/v5/middleware"
	"github.com/labstack/echo/v4"
)

var (
	mode       = flag.String("mode", "", "Run mode: v1, v2, or hybrid (default: auto-detect)")
	dataFolder = flag.String("data-folder", "", "Data folder path")
	migrate    = flag.Bool("migrate", false, "Run data migration from v1 to v2")
)

func main() {
	flag.Parse()

	// Load feature flags
	flags := integration.LoadFeatureFlags()

	log.Printf("Starting MMP Agent in %s mode", getMode(flags))
	log.Printf("Feature flags: %+v", flags)

	// Start pprof server
	go func() {
		http.ListenAndServe("localhost:8080", nil)
	}()

	// Determine run mode
	runMode := *mode
	if runMode == "" {
		runMode = flags.GetMode()
	}

	switch runMode {
	case "v1", "v1-only":
		runV1Only()
	case "v2", "v2-full":
		runV2Only()
	case "hybrid":
		runHybrid(flags)
	default:
		// Auto-detect based on flags
		if flags.IsV2Enabled() {
			runHybrid(flags)
		} else {
			runV1Only()
		}
	}
}

// runV1Only runs only the v1 system (legacy mode)
func runV1Only() {
	log.Println("Running in V1-only mode (legacy)")
	v1stlib.Run()
}

// runV2Only runs only the v2 system (modern mode)
func runV2Only() {
	log.Println("Running in V2-only mode (modern)")

	// Initialize v2 config
	dataFolderPath := *dataFolder
	if dataFolderPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Error getting user home directory: %v", err)
		}
		dataFolderPath = fmt.Sprintf("%s/.mmp", home)
	}

	if err := v2config.Load(dataFolderPath); err != nil {
		log.Fatalf("Error loading v2 config: %v", err)
	}

	// Initialize v2 database
	if err := v2database.Init(); err != nil {
		log.Fatalf("Error initializing v2 database: %v", err)
	}

	// Initialize v2 library
	lib, libHandler, _, err := v2library.New()
	if err != nil {
		log.Fatalf("Error initializing v2 library: %v", err)
	}

	// Create router
	r := chi.NewRouter()
	r.Use(cmw.Logger)
	r.Use(cmw.Recover())

	// Mount v2 API
	r.Route("/api/v2", func(r chi.Router) {
		r.Mount("/lib", libHandler)
	})

	// Start async scanning
	lib.ScanAsync(context.Background())

	// Start server
	port := v2config.Cfg.Server.Port
	log.Printf("Starting v2 server on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}

// runHybrid runs both v1 and v2 systems together
func runHybrid(flags *integration.FeatureFlags) {
	log.Println("Running in hybrid mode (v1 + v2)")

	// Check if migration is requested
	if *migrate {
		runMigration()
		return
	}

	// Create API router with versioning
	apiRouter := integration.NewAPIRouter(flags)

	// Setup v1 routes
	apiRouter.SetupV1Routes(func(e *echo.Echo) {
		// This would call the existing v1 route setup
		// For now, we'll register a simple health check
		e.GET("/health", func(c echo.Context) error {
			return c.JSON(200, map[string]string{"status": "ok", "version": "v1"})
		})
	})

	// Setup v2 routes if enabled
	if flags.EnableV2API {
		// Initialize v2 config
		dataFolderPath := *dataFolder
		if dataFolderPath == "" {
			home, _ := os.UserHomeDir()
			dataFolderPath = fmt.Sprintf("%s/.mmp", home)
		}

		if err := v2config.Load(dataFolderPath); err != nil {
			log.Printf("Warning: Could not load v2 config: %v", err)
		} else {
			// Initialize v2 database
			if err := v2database.Init(); err != nil {
				log.Printf("Warning: Could not initialize v2 database: %v", err)
			} else {
				// Initialize v2 library
				lib, libHandler, _, err := v2library.New()
				if err != nil {
					log.Printf("Warning: Could not initialize v2 library: %v", err)
				} else {
					apiRouter.SetupV2Routes(func(r chi.Router) {
						r.Mount("/lib", libHandler)
					})

					// Start async scanning
					lib.ScanAsync(context.Background())
				}
			}
		}
	}

	// Get unified handler
	handler := apiRouter.GetUnifiedHandler()

	// Start server
	port := 8000 // Default port, should come from config
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}

	// Graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	log.Printf("Starting hybrid server on port %d", port)
	log.Printf("API endpoints:")
	log.Printf("  - /api/v1/* (v1 routes)")
	if flags.EnableV2API {
		log.Printf("  - /api/v2/* (v2 routes)")
	}
	log.Printf("  - /api/* (legacy routes → v1)")

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

// runMigration performs data migration from v1 to v2
func runMigration() {
	log.Println("Running data migration from v1 to v2")
	log.Println("Migration not yet implemented - see core/integration/data_migration.go")
	// TODO: Implement actual migration logic
	os.Exit(0)
}

// getMode returns a friendly mode name
func getMode(flags *integration.FeatureFlags) string {
	return flags.GetMode()
}
