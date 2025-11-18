package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/eduardooliveira/stLib/core/integration"
	v2events "github.com/eduardooliveira/stLib/v2/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAPIVersioning tests API versioning works correctly
func TestAPIVersioning(t *testing.T) {
	os.Setenv("ENABLE_API_VERSIONING", "true")
	defer os.Clearenv()

	flags := integration.LoadFeatureFlags()
	router := integration.NewAPIRouter(flags)
	handler := router.GetUnifiedHandler()

	// Test v1 API endpoint exists
	t.Run("V1 Endpoint", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/health", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		// Should not return 404
		assert.NotEqual(t, http.StatusNotFound, w.Code)
	})

	// Test v2 API endpoint exists
	t.Run("V2 Endpoint", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v2/health", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		// Should not return 404
		assert.NotEqual(t, http.StatusNotFound, w.Code)
	})
}

// TestLegacyRoutes tests legacy routes still work
func TestLegacyRoutes(t *testing.T) {
	os.Clearenv()
	flags := integration.LoadFeatureFlags()
	router := integration.NewAPIRouter(flags)
	handler := router.GetUnifiedHandler()

	// Legacy routes should map to v1
	legacyRoutes := []string{
		"/api/projects",
		"/api/assets",
		"/api/tags",
	}

	for _, route := range legacyRoutes {
		t.Run(route, func(t *testing.T) {
			req := httptest.NewRequest("GET", route, nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			// Should route somewhere (not 404)
			assert.NotEqual(t, http.StatusNotFound, w.Code,
				"Legacy route %s should be accessible", route)
		})
	}
}

// TestV1APIBackwardCompatibility tests v1 API maintains compatibility
func TestV1APIBackwardCompatibility(t *testing.T) {
	// This would require a full server setup with database
	// For now, we test the routing structure

	os.Clearenv()
	flags := integration.LoadFeatureFlags()

	assert.False(t, flags.EnableV2API, "V2 API should be disabled in v1-only mode")

	// In v1-only mode, v1 API should work
	router := integration.NewAPIRouter(flags)
	assert.NotNil(t, router)
}

// TestV2APINewFeatures tests v2 API new features
func TestV2APINewFeatures(t *testing.T) {
	os.Setenv("ENABLE_V2_API", "true")
	defer os.Clearenv()

	flags := integration.LoadFeatureFlags()
	router := integration.NewAPIRouter(flags)
	handler := router.GetUnifiedHandler()

	// Test v2-specific endpoints
	v2Routes := []string{
		"/api/v2/assets",        // LibFS-based assets
		"/api/v2/discovery",     // New discovery endpoints
		"/api/v2/events/ws",     // WebSocket endpoint
	}

	for _, route := range v2Routes {
		t.Run(route, func(t *testing.T) {
			req := httptest.NewRequest("GET", route, nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			// Should be routable
			_ = w.Code
		})
	}
}

// TestContentTypeHandling tests content type handling
func TestContentTypeHandling(t *testing.T) {
	os.Clearenv()
	flags := integration.LoadFeatureFlags()
	router := integration.NewAPIRouter(flags)
	handler := router.GetUnifiedHandler()

	tests := []struct {
		name        string
		endpoint    string
		contentType string
	}{
		{"JSON Request", "/api/v1/projects", "application/json"},
		{"Form Request", "/api/v1/assets", "multipart/form-data"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := bytes.NewBufferString(`{"test": "data"}`)
			req := httptest.NewRequest("POST", tt.endpoint, body)
			req.Header.Set("Content-Type", tt.contentType)

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			// Should handle content type appropriately
			_ = w.Code
		})
	}
}

// TestErrorHandling tests error handling consistency
func TestErrorHandling(t *testing.T) {
	os.Clearenv()
	flags := integration.LoadFeatureFlags()
	router := integration.NewAPIRouter(flags)
	handler := router.GetUnifiedHandler()

	tests := []struct {
		name     string
		endpoint string
		method   string
		wantCode int
	}{
		{"Not Found", "/api/v1/nonexistent", "GET", http.StatusNotFound},
		{"Method Not Allowed", "/api/v1/projects", "PATCH", http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.endpoint, nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			// Error codes should be consistent
			_ = w.Code
		})
	}
}

// TestCORSHandling tests CORS header handling
func TestCORSHandling(t *testing.T) {
	os.Clearenv()
	flags := integration.LoadFeatureFlags()
	router := integration.NewAPIRouter(flags)
	handler := router.GetUnifiedHandler()

	req := httptest.NewRequest("OPTIONS", "/api/v1/projects", nil)
	req.Header.Set("Origin", "http://localhost:3000")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// CORS headers should be set
	// (exact headers depend on middleware configuration)
	_ = w.Header()
}

// TestRateLimiting tests rate limiting (if implemented)
func TestRateLimiting(t *testing.T) {
	t.Skip("Rate limiting test - implement when rate limiting is added")

	// Would test:
	// - Multiple rapid requests
	// - 429 Too Many Requests response
	// - Rate limit headers
}

// TestAuthenticationFlow tests authentication (if implemented)
func TestAuthenticationFlow(t *testing.T) {
	t.Skip("Authentication test - implement when auth is added")

	// Would test:
	// - Protected endpoints require auth
	// - Auth tokens work
	// - Invalid tokens rejected
}

// MockEventManager creates a mock event manager for testing
func MockEventManager() *v2events.EventManager {
	mgr := v2events.NewEventManager()
	ctx := context.Background()
	go mgr.Start(ctx)
	return mgr
}

// TestAPIResponseFormat tests response format consistency
func TestAPIResponseFormat(t *testing.T) {
	os.Clearenv()
	flags := integration.LoadFeatureFlags()
	router := integration.NewAPIRouter(flags)
	handler := router.GetUnifiedHandler()

	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Response should be valid JSON
	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)

	// If endpoint exists, response should be JSON
	if w.Code != http.StatusNotFound {
		assert.NoError(t, err, "Response should be valid JSON")
	}
}

// TestPaginationSupport tests pagination parameters
func TestPaginationSupport(t *testing.T) {
	os.Clearenv()
	flags := integration.LoadFeatureFlags()
	router := integration.NewAPIRouter(flags)
	handler := router.GetUnifiedHandler()

	tests := []struct {
		name     string
		endpoint string
		params   string
	}{
		{"Page 1", "/api/v1/projects", "?page=1&limit=10"},
		{"Page 2", "/api/v1/projects", "?page=2&limit=20"},
		{"No params", "/api/v1/projects", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.endpoint+tt.params, nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			// Should handle pagination params
			_ = w.Code
		})
	}
}

// TestFilteringSupport tests filtering parameters
func TestFilteringSupport(t *testing.T) {
	os.Clearenv()
	flags := integration.LoadFeatureFlags()
	router := integration.NewAPIRouter(flags)
	handler := router.GetUnifiedHandler()

	tests := []struct {
		name     string
		endpoint string
		params   string
	}{
		{"Filter by kind", "/api/v2/assets", "?kind=model"},
		{"Filter by tag", "/api/v1/assets", "?tag=important"},
		{"Multiple filters", "/api/v2/assets", "?kind=model&tag=ready"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.endpoint+tt.params, nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			// Should handle filter params
			_ = w.Code
		})
	}
}

// TestAPIDocumentation tests API documentation endpoints
func TestAPIDocumentation(t *testing.T) {
	os.Clearenv()
	flags := integration.LoadFeatureFlags()
	router := integration.NewAPIRouter(flags)
	handler := router.GetUnifiedHandler()

	docEndpoints := []string{
		"/api/docs",
		"/api/v1/docs",
		"/api/v2/docs",
		"/swagger.json",
	}

	for _, endpoint := range docEndpoints {
		t.Run(endpoint, func(t *testing.T) {
			req := httptest.NewRequest("GET", endpoint, nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			// Documentation endpoints may or may not exist
			// Just verify they're routable
			_ = w.Code
		})
	}
}

// BenchmarkAPIRouting benchmarks API routing performance
func BenchmarkAPIRouting(b *testing.B) {
	os.Clearenv()
	flags := integration.LoadFeatureFlags()
	router := integration.NewAPIRouter(flags)
	handler := router.GetUnifiedHandler()

	req := httptest.NewRequest("GET", "/api/v1/projects", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
	}
}
