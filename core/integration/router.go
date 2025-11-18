package integration

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// APIRouter manages versioned API routing
type APIRouter struct {
	// V1 uses Echo
	v1Echo *echo.Echo

	// V2 uses Chi
	v2Chi chi.Router

	// Flags for feature toggling
	flags *FeatureFlags
}

// NewAPIRouter creates a new API router with versioning support
func NewAPIRouter(flags *FeatureFlags) *APIRouter {
	return &APIRouter{
		v1Echo: echo.New(),
		v2Chi:  chi.NewRouter(),
		flags:  flags,
	}
}

// SetupV1Routes configures v1 API routes (Echo-based)
func (r *APIRouter) SetupV1Routes(setupFunc func(*echo.Echo)) {
	r.v1Echo.Use(middleware.Logger())
	r.v1Echo.Use(middleware.Recover())

	setupFunc(r.v1Echo)
}

// SetupV2Routes configures v2 API routes (Chi-based)
func (r *APIRouter) SetupV2Routes(setupFunc func(chi.Router)) {
	setupFunc(r.v2Chi)
}

// GetUnifiedHandler returns a single HTTP handler that routes to both v1 and v2
func (r *APIRouter) GetUnifiedHandler() http.Handler {
	router := chi.NewRouter()

	// V1 API routes
	router.Mount("/api/v1", echo.WrapHandler(r.v1Echo))

	// Legacy routes (map to v1 for backward compatibility)
	router.Mount("/api/projects", echo.WrapHandler(r.v1Echo.Group("/api/projects")))
	router.Mount("/api/tags", echo.WrapHandler(r.v1Echo.Group("/api/tags")))
	router.Mount("/api/assetTypes", echo.WrapHandler(r.v1Echo.Group("/api/assetTypes")))
	router.Mount("/api/tempfiles", echo.WrapHandler(r.v1Echo.Group("/api/tempfiles")))
	router.Mount("/api/system", echo.WrapHandler(r.v1Echo.Group("/api/system")))

	// V2 API routes (if enabled)
	if r.flags.EnableV2API || r.flags.EnableAPIVersioning {
		router.Mount("/api/v2", r.v2Chi)
	}

	// WebSocket events
	router.Mount("/api/events", echo.WrapHandler(r.v1Echo.Group("/api/events")))

	return router
}

// GetV1Echo returns the v1 Echo instance for configuration
func (r *APIRouter) GetV1Echo() *echo.Echo {
	return r.v1Echo
}

// GetV2Chi returns the v2 Chi router for configuration
func (r *APIRouter) GetV2Chi() chi.Router {
	return r.v2Chi
}
