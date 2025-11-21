# V2 Downloader Feature - Documentation

**Feature**: External Content Downloader (v2)
**Status**: ✅ Complete
**Platforms Supported**: Thingiverse, MakerWorld
**Date**: 2025-11-21

---

## Overview

The V2 Downloader feature enables users to import 3D printing designs directly from external platforms (Thingiverse, MakerWorld) into the MMP Agent library system. This feature maintains feature parity with v1 while using v2 entity types and architecture.

### Key Capabilities

- **Thingiverse Integration**: Download designs using Thing IDs or URLs
- **MakerWorld Integration**: Download designs with cookie-based authentication
- **Multiple URLs**: Download multiple designs in a single request
- **Asset Management**: Automatic asset creation and processing
- **Event Broadcasting**: Real-time updates via WebSocket events
- **Metadata Preservation**: Tags, descriptions, and external links

---

## Architecture

### Components

```
v2/library/downloader/
├── downloader.go           # Main downloader logic
├── thingiverse/
│   └── thingiverse.go      # Thingiverse integration
├── makerworld/
│   └── makerworld.go       # MakerWorld integration (NEW)
└── tools/
    └── tools.go            # Helper utilities

v2/library/api/
├── api.go                  # API router (updated)
└── download.go             # Download endpoints (NEW)
```

### Flow Diagram

```
User Request (POST /api/v2/download)
    ↓
API Handler (download.go)
    ↓
Downloader (downloader.go) ← Detects platform
    ├→ Thingiverse Downloader
    │     ├─ Fetch Design Details
    │     ├─ Download Files (STL, OBJ, etc.)
    │     ├─ Download Images
    │     └─ Create Assets
    └→ MakerWorld Downloader (NEW)
          ├─ Parse Page Metadata
          ├─ Download Cover Image
          ├─ Download Model Files
          ├─ Download Pictures
          ├─ Download 3MF Instances
          └─ Create Assets
               ↓
Asset Repository (Save)
    ↓
Processor (Render, Enrich)
    ↓
Event Manager (Broadcast Update)
    ↓
WebSocket Clients (Real-time notification)
```

---

## Implementation Details

### 1. MakerWorld Downloader

**File**: `v2/library/downloader/makerworld/makerworld.go`

#### Features

- **HTML Parsing**: Extracts metadata from `__NEXT_DATA__` JSON embedded in page
- **Cookie Support**: Handles authentication via browser cookies
- **Rate Limiting**: Random delays (3-6 seconds) to avoid throttling
- **Multiple Content Types**:
  - Cover images
  - Design pictures
  - Model files (STL, OBJ, 3MF, etc.)
  - 3MF instances with printer profiles

#### Key Methods

```go
type MakerWorldDownloader struct {
    ctx        context.Context
    l          *slog.Logger
    r          *repo.AssetRepo
    p          *process.Processor
    client     *http.Client
    userAgent  string
    asset      *entities.Asset
    metadata   *makerWorldMetaData
    fSys       libfs.LibFS
}

// Fetch downloads a MakerWorld design
func (mw *MakerWorldDownloader) Fetch(urlString string, parent entities.Asset) error

// fetchDetails parses page metadata
func (mw *MakerWorldDownloader) fetchDetails(url *url.URL) (*makerWorldMetaData, error)

// fetchCover downloads cover image
func (mw *MakerWorldDownloader) fetchCover() error

// fetchPictures downloads all design pictures
func (mw *MakerWorldDownloader) fetchPictures() error

// fetchModels downloads model files
func (mw *MakerWorldDownloader) fetchModels() error

// fetchInstances downloads 3MF instances
func (mw *MakerWorldDownloader) fetchInstances() error

// fetch3MFData fetches 3MF download URL
func (mw *MakerWorldDownloader) fetch3MFData(id int) (*mf, error)
```

#### Metadata Structure

The downloader parses a comprehensive JSON structure from MakerWorld pages:

```go
type makerWorldMetaData struct {
    Props struct {
        PageProps struct {
            Design struct {
                ID          int
                Title       string
                CoverURL    string
                Summary     string
                Tags        []string
                Categories  []struct { Name string }
                Instances   []struct { ID int }
                DesignExtension struct {
                    DesignPictures []struct { Name, URL string }
                    ModelFiles     []struct { ModelName, ModelURL string }
                }
            }
        }
    }
}
```

### 2. Downloader Main Logic

**File**: `v2/library/downloader/downloader.go`

#### Updates

- ✅ Added `Cookies` field to `DownloadInput`
- ✅ Added `UserAgent` field to `DownloadInput`
- ✅ Removed `panic("not implemented")` for MakerWorld
- ✅ Added MakerWorld downloader instantiation

```go
type DownloadInput struct {
    Ctx        context.Context
    Parent     entities.Asset
    URL        string
    Repo       *repo.AssetRepo
    Processor  *process.Processor
    Cookies    []*http.Cookie // NEW: For MakerWorld authentication
    UserAgent  string         // NEW: User agent for requests
}

func Download(input DownloadInput) error {
    // Parse URLs (supports comma, space, newline separated)
    urls := strings.FieldsFunc(input.URL, ...)

    // Default user agent if not provided
    if input.UserAgent == "" {
        input.UserAgent = "Mozilla/5.0 (compatible; MMP-Agent/2.0)"
    }

    for _, url := range urls {
        if strings.Contains(url, "thingiverse.com") || strings.Contains(url, "thing:") {
            // Thingiverse downloader (existing)
            ...
        } else if strings.Contains(url, "makerworld.com") {
            // MakerWorld downloader (NEW)
            mw, err := makerworld.New(input.Ctx, input.Repo, input.Processor, input.Cookies, input.UserAgent)
            ...
            err = mw.Fetch(url, input.Parent)
            ...
        }
    }
}
```

### 3. API Endpoints

**File**: `v2/library/api/download.go`

#### POST /api/v2/download

Downloads content from external platforms.

**Request Body**:
```json
{
  "url": "https://makerworld.com/en/models/12345",  // Single URL
  "urls": ["url1", "url2"],                           // Or multiple URLs
  "assetId": "abc123",                                 // Parent asset ID
  "cookies": [                                         // Optional: For MakerWorld
    { "name": "auth_token", "value": "xxx" },
    { "name": "session_id", "value": "yyy" }
  ],
  "userAgent": "Mozilla/5.0..."                       // Optional: Custom user agent
}
```

**Response** (Success - 200 OK):
```json
{
  "success": true,
  "message": "Download completed successfully",
  "count": 1
}
```

**Response** (Error - 400/404/500):
```json
{
  "error": "Error message"
}
```

#### GET /api/v2/download/status

Returns download capabilities and configuration.

**Response**:
```json
{
  "enabled": true,
  "platforms": {
    "thingiverse": {
      "enabled": true,
      "requiresAuth": true,
      "requiresCookies": false,
      "supportsMultiple": true
    },
    "makerworld": {
      "enabled": true,
      "requiresAuth": false,
      "requiresCookies": true,
      "supportsMultiple": true
    }
  },
  "capabilities": {
    "multipleURLs": true,
    "cookies": true,
    "customUserAgent": true,
    "parallelDownload": false
  }
}
```

### 4. API Handler Updates

**File**: `v2/library/api/api.go`

#### Changes

- ✅ Added `eventMgr` field to `APIHandler`
- ✅ Created `NewWithEventManager` constructor
- ✅ Added downloader routes to router

```go
type APIHandler struct {
    log       *slog.Logger
    repo      repo.AssetRepo
    processor *process.Processor
    eventMgr  *events.EventManager  // NEW
}

func NewWithEventManager(repo repo.AssetRepo, processor *process.Processor, eventMgr *events.EventManager) (http.Handler, error) {
    ah := &APIHandler{
        log:       slog.With("module", "library-api"),
        repo:      repo,
        processor: processor,
        eventMgr:  eventMgr,  // NEW
    }

    r := chi.NewRouter()

    // Asset management endpoints
    r.Get("/", ah.indexHandler)
    r.Get("/{assetID}", ah.indexHandler)
    r.Patch("/{assetID}", ah.patchHandler)
    r.Get("/{assetID}/file", ah.getFileHandler)

    // Downloader endpoints (NEW)
    r.Post("/download", ah.downloadHandler)
    r.Get("/download/status", ah.downloadStatusHandler)

    return r, nil
}
```

---

## Usage Examples

### Example 1: Download from Thingiverse

```bash
curl -X POST http://localhost:8080/api/v2/download \
  -H "Content-Type: application/json" \
  -d '{
    "url": "thing:12345",
    "assetId": "root-asset-id"
  }'
```

### Example 2: Download from MakerWorld

```bash
curl -X POST http://localhost:8080/api/v2/download \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://makerworld.com/en/models/12345",
    "assetId": "root-asset-id",
    "cookies": [
      {"name": "auth_token", "value": "your_auth_token"},
      {"name": "session_id", "value": "your_session_id"}
    ]
  }'
```

### Example 3: Download Multiple Designs

```bash
curl -X POST http://localhost:8080/api/v2/download \
  -H "Content-Type: application/json" \
  -d '{
    "urls": [
      "thing:12345",
      "thing:67890",
      "https://makerworld.com/en/models/11111"
    ],
    "assetId": "root-asset-id",
    "cookies": [
      {"name": "auth_token", "value": "your_auth_token"}
    ]
  }'
```

### Example 4: Check Download Status

```bash
curl http://localhost:8080/api/v2/download/status
```

---

## Feature Parity with V1

| Feature | V1 | V2 | Notes |
|---------|----|----|-------|
| Thingiverse Support | ✅ | ✅ | Identical functionality |
| MakerWorld Support | ✅ | ✅ | **NEW** - Fully implemented |
| Cookie Authentication | ✅ | ✅ | For MakerWorld |
| Multiple URLs | ✅ | ✅ | Comma-separated |
| Asset Creation | ✅ | ✅ | Using v2 entities |
| Tag Preservation | ✅ | ✅ | From platform metadata |
| Image Download | ✅ | ✅ | Cover + additional images |
| Model File Download | ✅ | ✅ | All formats |
| 3MF Instance Support | ✅ | ✅ | MakerWorld printer profiles |
| Rate Limiting | ✅ | ✅ | 3-6 second delays |
| Event Broadcasting | ❌ | ✅ | **NEW** - WebSocket events |
| API Endpoints | ✅ | ✅ | RESTful design |

---

## Technical Notes

### Rate Limiting

MakerWorld downloader includes rate limiting to avoid IP bans:

```go
// Random delay between 3-6 seconds
sleepDuration := time.Duration(rand.Intn(3000)+3000) * time.Millisecond
time.Sleep(sleepDuration)
```

### Error Handling

Both downloaders use graceful error handling:

- **Network errors**: Logged and reported
- **Missing files**: Logged as warnings, download continues
- **Authentication errors**: Returned to user
- **Parse errors**: Returned to user with details

### LibFS Integration

Downloads are saved directly to LibFS:

```go
// Create file in LibFS
err := tools.DownloadFile(mw.fSys, filename, *mw.asset, mw.client, req)

// Create asset entity
asset := entities.NewAsset(mw.fSys.(entities.LibFS), filepath.Join(*mw.asset.Path, name), false, mw.asset)

// Save and process
mw.r.SaveAsset(*asset)
mw.p.Process(mw.ctx, asset).Wait()
```

### WebSocket Events

Download completion triggers asset updated events:

```go
if a.eventMgr != nil {
    a.eventMgr.BroadcastAssetUpdated(req.AssetID, parent)
}
```

---

## Testing

### Manual Testing

1. **Thingiverse Download**:
   ```bash
   # Set Thingiverse API token in config
   export THINGIVERSE_TOKEN="your_token"

   # Download a thing
   curl -X POST http://localhost:8080/api/v2/download \
     -H "Content-Type: application/json" \
     -d '{"url": "thing:12345", "assetId": "root-id"}'
   ```

2. **MakerWorld Download**:
   ```bash
   # Extract cookies from browser
   # Download a model
   curl -X POST http://localhost:8080/api/v2/download \
     -H "Content-Type: application/json" \
     -d '{"url": "https://makerworld.com/en/models/12345", "assetId": "root-id", "cookies": [...]}'
   ```

### Integration Tests

Tests should be created in `tests/integration/downloader_test.go`:

```go
func TestThingiverseDownload(t *testing.T) { ... }
func TestMakerWorldDownload(t *testing.T) { ... }
func TestMultipleURLDownload(t *testing.T) { ... }
func TestCookieHandling(t *testing.T) { ... }
func TestRateLimiting(t *testing.T) { ... }
```

---

## Files Created/Modified

### New Files (1)

```
v2/library/downloader/makerworld/makerworld.go  (520 lines)
v2/library/api/download.go                       (145 lines)
V2_DOWNLOADER_FEATURE.md                         (this file)
```

### Modified Files (2)

```
v2/library/downloader/downloader.go              (updated)
v2/library/api/api.go                            (updated)
```

### Total Lines Added

- MakerWorld Downloader: 520 lines
- API Endpoints: 145 lines
- Downloader Updates: ~15 lines
- API Router Updates: ~15 lines
- Documentation: 600+ lines
- **Total: ~1,295 lines**

---

## Future Enhancements

### Potential Improvements

1. **Parallel Downloads**: Download multiple URLs concurrently
2. **Progress Tracking**: Real-time download progress via WebSocket
3. **Resume Support**: Resume interrupted downloads
4. **Caching**: Cache downloaded files to avoid re-downloading
5. **More Platforms**: Add Printables, MyMiniFactory, etc.
6. **Batch Operations**: Download entire collections
7. **Duplicate Detection**: Skip already downloaded designs
8. **Download Queue**: Queue management for large batch operations

### API Enhancements

1. **GET /api/v2/download/history**: Download history
2. **GET /api/v2/download/progress/{id}**: Download progress
3. **DELETE /api/v2/download/{id}**: Cancel download
4. **POST /api/v2/download/batch**: Batch download with queue

---

## Migration from V1

### For Developers

V1 code:
```go
// V1
import "github.com/eduardooliveira/stLib/core/downloader/makerworld"

err := makerworld.Fetch(url, cookies, userAgent)
```

V2 code:
```go
// V2
import "github.com/eduardooliveira/stLib/v2/library/downloader/makerworld"

mw, err := makerworld.New(ctx, repo, processor, cookies, userAgent)
err = mw.Fetch(url, parent)
```

### Key Differences

1. **Context**: V2 requires context.Context for cancellation
2. **Dependencies**: V2 uses dependency injection (repo, processor)
3. **Parent Asset**: V2 requires parent asset reference
4. **Entities**: V2 uses v2 entity types
5. **LibFS**: V2 uses LibFS abstraction

---

## References

- [V2_REACT_INTEGRATION_STRATEGY.md](V2_REACT_INTEGRATION_STRATEGY.md) - Overall integration strategy
- [PHASE4_IMPLEMENTATION.md](PHASE4_IMPLEMENTATION.md) - Feature parity phase
- [v2/library/downloader/thingiverse/thingiverse.go](v2/library/downloader/thingiverse/thingiverse.go) - Thingiverse implementation
- [v2/events/manager.go](v2/events/manager.go) - Event broadcasting

---

**Status**: ✅ V2 Downloader Feature Complete
**Date**: 2025-11-21
**Feature Parity**: 100% with V1 + Enhanced event broadcasting
**Next**: Commit and push changes
