package api

import (
	"encoding/json"
	"net/http"

	"github.com/eduardooliveira/stLib/v2/library/downloader"
	"github.com/eduardooliveira/stLib/v2/library/repo"
)

// downloadHandler handles POST /api/v2/download
// Downloads content from external sources (Thingiverse, MakerWorld)
func (a *APIHandler) downloadHandler(w http.ResponseWriter, r *http.Request) {
	a.log.Info("Download request received")

	// Parse request body
	var req struct {
		URL       string   `json:"url"`       // Single URL
		URLs      []string `json:"urls"`      // Multiple URLs
		AssetID   string   `json:"assetId"`   // Parent asset ID
		Cookies   []Cookie `json:"cookies"`   // For MakerWorld authentication
		UserAgent string   `json:"userAgent"` // Optional user agent
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.log.Error("Failed to decode request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Combine single URL and URLs array
	urls := req.URLs
	if req.URL != "" {
		urls = append(urls, req.URL)
	}

	if len(urls) == 0 {
		a.log.Warn("No URLs provided")
		http.Error(w, "No URLs provided", http.StatusBadRequest)
		return
	}

	// Get parent asset
	if req.AssetID == "" {
		a.log.Warn("No asset ID provided")
		http.Error(w, "Asset ID is required", http.StatusBadRequest)
		return
	}

	parent, err := a.repo.GetAsset(req.AssetID)
	if err != nil {
		a.log.Error("Failed to get parent asset", "error", err, "assetId", req.AssetID)
		http.Error(w, "Parent asset not found", http.StatusNotFound)
		return
	}

	// Convert cookies
	var cookies []*http.Cookie
	for _, c := range req.Cookies {
		cookies = append(cookies, &http.Cookie{
			Name:  c.Name,
			Value: c.Value,
			Path:  "/",
		})
	}

	// User agent
	userAgent := req.UserAgent
	if userAgent == "" {
		userAgent = r.UserAgent()
	}

	// Get concrete repo type
	assetRepo, ok := a.repo.(*repo.AssetRepo)
	if !ok {
		a.log.Error("Invalid repository type")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Download each URL
	for _, url := range urls {
		a.log.Info("Downloading", "url", url)

		input := downloader.DownloadInput{
			Ctx:       r.Context(),
			Parent:    *parent,
			URL:       url,
			Repo:      assetRepo,
			Processor: a.processor,
			Cookies:   cookies,
			UserAgent: userAgent,
		}

		if err := downloader.Download(input); err != nil {
			a.log.Error("Download failed", "url", url, "error", err)
			http.Error(w, "Download failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Broadcast event if event manager is available
	if a.eventMgr != nil {
		a.eventMgr.BroadcastAssetUpdated(req.AssetID, parent)
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Download completed successfully",
		"count":   len(urls),
	})
}

// downloadStatusHandler handles GET /api/v2/download/status
// Returns download capabilities and configuration
func (a *APIHandler) downloadStatusHandler(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"enabled": true,
		"platforms": map[string]interface{}{
			"thingiverse": map[string]interface{}{
				"enabled":          true,
				"requiresAuth":     true,
				"requiresCookies":  false,
				"supportsMultiple": true,
			},
			"makerworld": map[string]interface{}{
				"enabled":          true,
				"requiresAuth":     false,
				"requiresCookies":  true,
				"supportsMultiple": true,
			},
		},
		"capabilities": map[string]bool{
			"multipleURLs":   true,
			"cookies":        true,
			"customUserAgent": true,
			"parallelDownload": false,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// Cookie represents a browser cookie for authentication
type Cookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
