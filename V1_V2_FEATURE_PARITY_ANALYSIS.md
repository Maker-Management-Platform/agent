# V1 to V2 Feature Parity Analysis

**Date**: 2025-11-21
**Status**: Gap Analysis Complete

---

## Executive Summary

This document provides a comprehensive analysis of feature parity between v1 and v2 architectures, identifying what features are present in v1 but missing in v2.

### Overall Status

| Category | V1 Features | V2 Features | Parity % |
|----------|-------------|-------------|----------|
| **Core API Endpoints** | ~35 endpoints | ~6 endpoints | **17%** |
| **Integrations** | Printers, Slicer | Printers (partial) | **30%** |
| **Data Management** | Full CRUD | Partial CRUD | **40%** |
| **Event System** | SSE | WebSocket | **80%** |
| **Downloader** | Thingiverse + MakerWorld | Thingiverse + MakerWorld | **100%** |
| **Discovery/Processing** | Full | Full | **100%** |

**Overall Feature Parity: ~45%**

---

## Missing Features by Category

### 1. Projects API (Missing in V2)

#### V1 Endpoints (All Missing in V2):
```
GET    /api/projects              - List all projects
GET    /api/projects/list         - List projects (alternate)
GET    /api/projects/:uuid        - Get project by UUID
GET    /api/projects/:uuid/discover - Discover assets in project
POST   /api/projects              - Create new project
POST   /api/projects/:uuid        - Update/save project
POST   /api/projects/:uuid/move   - Move project
POST   /api/projects/:uuid/image  - Set main image
POST   /api/projects/:uuid/delete - Delete project
```

#### Capabilities in V1:
- ✅ Create projects with file uploads
- ✅ Update project metadata (name, description, tags)
- ✅ Delete projects
- ✅ Move projects to different locations
- ✅ Set main/default image
- ✅ Discover assets within project
- ✅ List all projects
- ✅ Get project details

#### V2 Status:
- ❌ **No project-specific API endpoints**
- ⚠️ V2 uses root assets instead of projects
- ⚠️ Mapping: v1 "Project" → v2 "Root Asset"

---

### 2. Assets API (Partial in V2)

#### V1 Endpoints:
```
GET    /api/projects/:uuid/assets         - List assets in project
POST   /api/projects/:uuid/assets         - Create new asset
GET    /api/projects/:uuid/assets/:id     - Get asset by ID
GET    /api/projects/:uuid/assets/:id/file - Get asset file
POST   /api/projects/:uuid/assets/:id/delete - Delete asset
```

#### V2 Endpoints (Existing):
```
GET    /api/v2/                    - Index/list assets (partial)
GET    /api/v2/{assetID}           - Get asset by ID
PATCH  /api/v2/{assetID}           - Update asset
GET    /api/v2/{assetID}/file      - Get asset file
```

#### Missing in V2:
- ❌ **POST /api/v2/assets** - Create new asset endpoint
- ❌ **DELETE /api/v2/{assetID}** - Delete asset endpoint
- ❌ **Project-scoped asset listing** (v2 has flat asset listing)
- ❌ **Asset file upload** (create asset with file)

---

### 3. Tags API (Missing in V2)

#### V1 Endpoints:
```
GET    /api/tags                   - Get all tags
```

#### V2 Status:
- ❌ **No tag management endpoints**
- ⚠️ Tags exist in v2 entities but no API to manage them
- ⚠️ Tags can be added to assets but not listed separately

---

### 4. System API (Missing in V2)

#### V1 Endpoints:
```
GET    /api/system/paths           - Get unregistered paths
GET    /api/system/settings        - Get system settings
POST   /api/system/settings        - Save system settings
POST   /api/system/discovery       - Run full discovery
POST   /api/system/subscribe/:session    - Subscribe to system events
POST   /api/system/unsubscribe/:session  - Unsubscribe from system events
```

#### Capabilities in V1:
- ✅ List folders not yet registered as projects
- ✅ Get/update system configuration
- ✅ Manually trigger discovery
- ✅ Subscribe to system-level events (SSE)

#### V2 Status:
- ❌ **No system management endpoints**
- ❌ **No configuration management API**
- ❌ **No manual discovery trigger**
- ⚠️ V2 uses different event system (WebSocket vs SSE)

---

### 5. Temp Files API (Missing in V2)

#### V1 Endpoints:
```
GET    /api/tempfiles              - List temporary files
POST   /api/tempfiles/:uuid/delete - Delete temp file
```

#### Purpose in V1:
- Slicer integration uploads files to temp location
- Files are discovered and listed
- Can be deleted when no longer needed

#### V2 Status:
- ❌ **No temp file management**
- ⚠️ May not be needed if slicer integration is redesigned

---

### 6. Asset Types API (Missing in V2)

#### V1 Endpoints:
```
GET    /api/assettypes             - Get all asset types
```

#### Purpose in V1:
- Lists registered asset types (model, gcode, image, etc.)
- Frontend uses this to know what file types are supported

#### V2 Status:
- ❌ **No asset type listing endpoint**
- ⚠️ V2 uses `Kind` field in assets but no API to list types

---

### 7. Printer Integration API (Partially Implemented)

#### V1 Endpoints (Full CRUD + Operations):
```
GET    /api/printers                      - List all printers
GET    /api/printers/:uuid                - Get printer by UUID
POST   /api/printers                      - Create new printer
POST   /api/printers/:uuid                - Update printer
POST   /api/printers/:uuid/delete         - Delete printer
POST   /api/printers/:uuid/test           - Test connection
POST   /api/printers/:uuid/send/:id       - Send file to printer
GET    /api/printers/:uuid/stream         - Stream camera
GET    /api/printers/:uuid/status         - Get printer status
POST   /api/printers/:uuid/subscribe/:session   - Subscribe to printer events
POST   /api/printers/:uuid/unsubscribe/:session - Unsubscribe from printer events
```

#### V2 Status (Phase 4 - Partial):
- ✅ V2 Printer Manager exists (`v2/integrations/printers/manager.go`)
- ✅ OctoPrint adapter implemented
- ✅ Klipper adapter implemented
- ✅ State management
- ❌ **No REST API endpoints for printers**
- ❌ **No CRUD operations via API**
- ❌ **No camera streaming**
- ❌ **No file sending API**

---

### 8. Slicer Integration API (Missing in V2)

#### V1 Endpoints (Moonraker/Klipper Compatible):
```
GET    /server/info                - Server information
POST   /server/files/upload        - Upload file (slicer compatibility)
GET    /api/version                - API version
```

#### Purpose in V1:
- Allows slicers (PrusaSlicer, Cura, etc.) to upload GCode directly
- Moonraker API compatibility layer
- Files go to temp directory for review

#### V2 Status:
- ❌ **No slicer integration endpoints**
- ❌ **No Moonraker compatibility**
- ❌ **No temp file upload**

---

### 9. Events API (Different Implementation)

#### V1 System (SSE - Server-Sent Events):
```
GET    /api/events/:session        - SSE event stream
```

**Capabilities**:
- Subscribe to multiple event channels
- System events, printer events, etc.
- Client specifies session ID

#### V2 System (WebSocket):
```
(Handled via WebSocket connection, not REST endpoints)
```

**Capabilities**:
- ✅ WebSocket-based real-time events
- ✅ Event broadcasting to all clients
- ✅ Asset lifecycle events
- ✅ Printer status events

**Differences**:
- V1 uses SSE (one-way, HTTP-based)
- V2 uses WebSocket (two-way, persistent connection)
- Both achieve similar goals, different implementations
- V2 is more modern and efficient

---

## Detailed Feature Matrix

### Core CRUD Operations

| Feature | V1 | V2 | Gap |
|---------|----|----|-----|
| **Projects** |
| List projects | ✅ | ❌ | No project concept |
| Create project | ✅ | ❌ | Use root assets |
| Get project | ✅ | ❌ | Use root assets |
| Update project | ✅ | ❌ | Use PATCH asset |
| Delete project | ✅ | ❌ | Need DELETE endpoint |
| Move project | ✅ | ❌ | Not supported |
| Set project image | ✅ | ❌ | Not supported |
| **Assets** |
| List assets | ✅ | ✅ | ✓ |
| Create asset | ✅ | ❌ | Need POST endpoint |
| Get asset | ✅ | ✅ | ✓ |
| Update asset | ✅ | ✅ | ✓ |
| Delete asset | ✅ | ❌ | Need DELETE endpoint |
| Get asset file | ✅ | ✅ | ✓ |
| Upload asset file | ✅ | ❌ | Need multipart upload |
| **Tags** |
| List tags | ✅ | ❌ | No endpoint |
| **Printers** |
| List printers | ✅ | ❌ | Manager exists, no API |
| Create printer | ✅ | ❌ | Need POST endpoint |
| Update printer | ✅ | ❌ | Need PATCH endpoint |
| Delete printer | ✅ | ❌ | Need DELETE endpoint |
| Test connection | ✅ | ❌ | Need endpoint |
| Send file | ✅ | ❌ | Need endpoint |
| Stream camera | ✅ | ❌ | Need endpoint |
| Get status | ✅ | ❌ | Events exist, no direct endpoint |

### System Management

| Feature | V1 | V2 | Gap |
|---------|----|----|-----|
| Get settings | ✅ | ❌ | No config API |
| Save settings | ✅ | ❌ | No config API |
| List unregistered paths | ✅ | ❌ | Not needed? |
| Trigger discovery | ✅ | ❌ | No manual trigger |
| System events | ✅ | ✅ | Different (SSE vs WS) |

### Integrations

| Feature | V1 | V2 | Gap |
|---------|----|----|-----|
| **Downloaders** |
| Thingiverse | ✅ | ✅ | ✓ |
| MakerWorld | ✅ | ✅ | ✓ |
| **Printers** |
| OctoPrint | ✅ | ✅ | Backend only |
| Klipper | ✅ | ✅ | Backend only |
| Moonraker | ✅ | ❌ | Not implemented |
| **Slicer** |
| File upload | ✅ | ❌ | No slicer API |
| Moonraker compat | ✅ | ❌ | No slicer API |

---

## Priority Ranking

### P0 - Critical (Blocks Basic Functionality)

1. **Asset CRUD Endpoints**
   - POST /api/v2/assets - Create asset with file upload
   - DELETE /api/v2/{assetID} - Delete asset
   - **Impact**: Cannot create or delete assets via API

2. **Printer Management API**
   - Full CRUD for printers
   - Send file to printer
   - **Impact**: Printer manager exists but unusable from frontend

### P1 - High (Major Feature Gaps)

3. **Tags API**
   - GET /api/v2/tags - List all tags
   - **Impact**: Cannot list or filter by tags effectively

4. **System API**
   - GET /api/v2/system/settings
   - POST /api/v2/system/settings
   - POST /api/v2/system/discovery
   - **Impact**: Cannot configure system or trigger discovery

5. **Project API** (or Root Asset equivalent)
   - Need project-like operations on root assets
   - **Impact**: V2 has different paradigm but needs similar capabilities

### P2 - Medium (Nice to Have)

6. **Slicer Integration**
   - Moonraker-compatible endpoints
   - **Impact**: Cannot upload directly from slicer software

7. **Printer Camera Streaming**
   - GET /api/v2/printers/{id}/stream
   - **Impact**: Cannot view printer cameras

8. **Asset Types API**
   - GET /api/v2/asset-types
   - **Impact**: Frontend doesn't know supported types

### P3 - Low (Minor Features)

9. **Temp Files API**
   - May not be needed depending on slicer implementation

10. **Project Move**
    - Complex operation, may not be needed in v2

---

## Recommended Implementation Order

### Phase 1: Essential CRUD (Week 1-2)
1. Asset creation endpoint (POST /api/v2/assets)
2. Asset deletion endpoint (DELETE /api/v2/{assetID})
3. Tags listing endpoint (GET /api/v2/tags)

### Phase 2: Printer API (Week 2-3)
4. Printer CRUD endpoints (all 5 endpoints)
5. Send file to printer endpoint
6. Test connection endpoint

### Phase 3: System Management (Week 3-4)
7. System settings endpoints (GET/POST)
8. Manual discovery trigger
9. Asset types listing

### Phase 4: Advanced Features (Week 4-5)
10. Camera streaming endpoint
11. Slicer integration endpoints
12. Temp file management (if needed)

---

## Migration Considerations

### Conceptual Differences

**v1 Projects vs v2 Root Assets**:
- v1: Projects are top-level containers with explicit CRUD
- v2: Root assets are just assets with `NodeKind=root`
- **Impact**: Need to expose root asset operations as "project-like" API

**v1 SSE vs v2 WebSocket**:
- v1: Server-Sent Events for one-way communication
- v2: WebSocket for bi-directional communication
- **Impact**: Frontend needs WebSocket client instead of EventSource

### Data Model Differences

**IDs**:
- v1: Auto-increment integers + UUIDs
- v2: MD5-based string IDs
- **Impact**: API responses have different ID formats

**Relationships**:
- v1: Project → Assets (one-to-many)
- v2: Root Asset → Nested Assets (tree structure)
- **Impact**: API responses have different nesting

---

## Files That Need to Be Created

Based on this analysis, here are the files that need to be created for v2 feature parity:

### High Priority

```
v2/library/api/
├── create.go              # Asset creation with file upload
├── delete.go              # Asset deletion
└── tags.go                # Tag management

v2/integrations/printers/api/
├── api.go                 # Printer API router
├── crud.go                # Printer CRUD operations
├── operations.go          # Send file, test connection
└── streaming.go           # Camera streaming
```

### Medium Priority

```
v2/system/api/
├── api.go                 # System API router
├── settings.go            # Configuration management
├── discovery.go           # Manual discovery trigger
└── info.go                # System information

v2/library/api/
└── asset_types.go         # Asset type listing
```

### Lower Priority

```
v2/slicer/
├── api.go                 # Slicer integration router
├── upload.go              # File upload endpoint
└── moonraker.go           # Moonraker compatibility

v2/library/api/
└── temp_files.go          # Temp file management
```

---

## Conclusion

### Current State
- **V2 has ~45% feature parity with V1**
- Core functionality exists (discovery, processing, downloader)
- Major gaps in API endpoints and integrations

### To Reach 100% Parity
- Implement ~30 missing API endpoints
- Complete printer integration API
- Add system management endpoints
- Implement slicer compatibility layer

### Estimated Effort
- **4-5 weeks** of development
- **~3,000-4,000 lines of code**
- **15-20 new API endpoint files**

### Next Steps
1. Prioritize P0 endpoints (asset CRUD, printer API)
2. Implement in phases as outlined above
3. Test each phase thoroughly
4. Update frontend to use v2 APIs
5. Deprecate v1 endpoints gradually

---

**Document Status**: Complete
**Last Updated**: 2025-11-21
**Next Action**: Begin Phase 1 implementation (Essential CRUD)
