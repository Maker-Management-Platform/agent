# V2-React Integration Strategy for MMP Agent

**Document Version:** 1.0
**Date:** November 6, 2025
**Author:** Claude AI
**Status:** Draft for Review

---

## Executive Summary

This document outlines a comprehensive strategy for integrating the v2-react branch's modern filesystem structure into the v1 (main) version of the Maker Management Platform (MMP) Agent. The v2-react branch represents a significant architectural evolution featuring:

- **Decoupled React 18 frontend** (52 TypeScript files)
- **Modular v2 backend architecture** (60 Go files)
- **Virtual filesystem abstraction** supporting local, git, and bundle sources
- **Desktop application support** via Wails framework
- **Modern tooling** (Vite, TanStack Query, Mantine UI)

**Key Goal:** Migrate from the monolithic v1 architecture to the modern v2 structure while maintaining backward compatibility and minimizing service disruption.

---

## Table of Contents

1. [Architecture Comparison](#1-architecture-comparison)
2. [Key Differences](#2-key-differences)
3. [Integration Strategy Overview](#3-integration-strategy-overview)
4. [Phase 1: Foundation & Preparation](#4-phase-1-foundation--preparation)
5. [Phase 2: Backend Migration](#5-phase-2-backend-migration)
6. [Phase 3: Frontend Integration](#6-phase-3-frontend-integration)
7. [Phase 4: Feature Parity](#7-phase-4-feature-parity)
8. [Phase 5: Testing & Validation](#8-phase-5-testing--validation)
9. [Phase 6: Deployment & Rollout](#9-phase-6-deployment--rollout)
10. [Risks & Mitigation](#10-risks--mitigation)
11. [Success Metrics](#11-success-metrics)
12. [Rollback Plan](#12-rollback-plan)

---

## 1. Architecture Comparison

### 1.1 V1 (Main Branch) Architecture

```
┌─────────────────────────────────────────────────┐
│                   main.go                       │
│              (Entry Point + pprof)              │
└────────────────────┬────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────┐
│                  stlib.go                       │
│           (Monolithic Orchestrator)             │
│  - Database Init                                │
│  - Processing Pipeline                          │
│  - Echo HTTP Server                             │
└────────────────────┬────────────────────────────┘
                     │
         ┌───────────┴───────────┐
         ▼                       ▼
┌────────────────┐      ┌───────────────────┐
│   /core/       │      │   /frontend/      │
│  (75 Go files) │      │   (Empty)         │
│                │      │                   │
│  - api/        │      │  No UI code       │
│  - processing/ │      │  Likely external  │
│  - data/       │      │  or server-side   │
│  - entities/   │      └───────────────────┘
│  - events/     │
│  - integr...   │
└────────────────┘

Technology Stack:
- Backend: Go 1.22, Echo v4, GORM, SQLite
- Frontend: Unknown/External (frontend/ is empty)
- Config: config.toml (9 lines)
- Deployment: Docker, standalone binary
```

**Strengths:**
- ✅ Mature, stable codebase (5,437 LOC)
- ✅ Comprehensive feature set (projects, assets, tags, printers)
- ✅ Real-time events via WebSocket
- ✅ Production-ready integrations (OctoPrint, Klipper, Thingiverse)

**Weaknesses:**
- ❌ No modern frontend (empty directory)
- ❌ Monolithic structure (stlib.go orchestrates everything)
- ❌ Direct filesystem coupling
- ❌ Limited scalability/modularity
- ❌ No desktop application support

---

### 1.2 V2-React Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     /cmd/                                   │
│   ┌────────────┐  ┌────────────┐  ┌─────────────┐         │
│   │ command/   │  │ desktop/   │  │ desktopgen/ │         │
│   │ (Web)      │  │ (Wails)    │  │ (Codegen)   │         │
│   └──────┬─────┘  └──────┬─────┘  └─────────────┘         │
└──────────┼────────────────┼──────────────────────────────────┘
           │                │
           │    ┌───────────┴──────────────┐
           │    │                           │
           ▼    ▼                           │
┌────────────────────────┐                 │
│       /v2/             │                 │
│  (Clean Architecture)  │                 │
│                        │                 │
│  ┌──────────────────┐  │                 │
│  │ /library/        │  │                 │
│  │ (37 files)       │  │                 │
│  │  - api/          │  │                 │
│  │  - entities/     │  │                 │
│  │  - libfs/        │  │◄────────────────┤
│  │  - repo/         │  │                 │
│  │  - process/      │  │                 │
│  │  - discovery/    │  │                 │
│  │  - downloader/   │  │                 │
│  └──────────────────┘  │                 │
│                        │                 │
│  ┌──────────────────┐  │                 │
│  │ /config/         │  │                 │
│  │ /database/       │  │                 │
│  │ /api/            │  │                 │
│  │ /utils/          │  │                 │
│  │ /web/ (templ)    │  │                 │
│  └──────────────────┘  │                 │
└────────────────────────┘                 │
                                           │
┌──────────────────────────────────────────▼───┐
│            /frontend/                        │
│        (React 18 + TypeScript)               │
│                                              │
│  ┌─────────────────────────────────┐        │
│  │ /src/                            │        │
│  │  - App.tsx (Mantine theme)      │        │
│  │  - /lib/ (Asset management)     │        │
│  │    - pages/                     │        │
│  │    - components/                │        │
│  │    - fetchers/ (React Query)    │        │
│  │  - /dashboard/                  │        │
│  │  - /core/ (shared infra)        │        │
│  └─────────────────────────────────┘        │
│                                              │
│  Build: Vite, TypeScript 5.7                │
│  UI: Mantine v7, Tabler Icons               │
│  3D: Three.js                                │
└──────────────────────────────────────────────┘

Technology Stack:
- Backend: Go 1.22, Chi v5, GORM, SQLite/Postgres
- Frontend: React 18 + TypeScript, Mantine, Three.js
- Desktop: Wails v2
- Dev Tools: Air hot-reload, Vite HMR
- Config: Viper + config.toml (TOML array support)
```

**Strengths:**
- ✅ Modern React frontend with TypeScript
- ✅ Clean, modular architecture
- ✅ Virtual filesystem abstraction (local, git, bundles)
- ✅ Desktop app support via Wails
- ✅ Better separation of concerns
- ✅ Hot reload development (Air + Vite)
- ✅ Enterprise UI components (Mantine)

**Weaknesses:**
- ❌ Feature incomplete (documented as "under development")
- ❌ Missing integrations from v1 (OctoPrint, Klipper, printers)
- ❌ No WebSocket events (yet)
- ❌ Limited testing/validation
- ❌ README mentions "dropping React" but code shows React (documentation lag)

---

## 2. Key Differences

### 2.1 Structural Differences

| Aspect | V1 (Main) | V2-React |
|--------|-----------|----------|
| **Entry Point** | Single `main.go` → `stlib.Run()` | Multiple: `/cmd/command/`, `/cmd/desktop/` |
| **Backend Code** | `/core/` (75 Go files) | `/v2/` (44 Go files) + legacy `/core/` |
| **Frontend** | Empty `/frontend/` | Populated `/frontend/` (52 files) |
| **HTTP Framework** | Echo v4 | Chi v5 |
| **Router** | Echo handlers | Chi router + middleware |
| **Config Management** | Simple TOML parsing | Viper (advanced TOML) |
| **Database Layer** | `/core/data/database/` | `/v2/database/` |
| **File System** | Direct OS filesystem | Virtual FS abstraction (`LibFS` interface) |
| **Desktop Support** | None | Wails v2 integration |
| **Dev Tooling** | Manual rebuild | Air hot-reload + Vite HMR |

---

### 2.2 Data Model Evolution

#### V1 Asset Model (core/entities/asset.go)
```go
type Asset struct {
    ID           uint
    Name         string
    ProjectID    uint
    Project      *Project
    AssetTypeID  uint
    AssetType    *AssetType
    Path         string
    // ... other fields
}
```

#### V2 Asset Model (v2/library/entities/asset.go)
```go
type Asset struct {
    ID           string         // MD5 hash instead of auto-increment
    Label        *string        // Renamed from Name
    Description  *string        // New field
    Path         *string        // Pointer for optional
    Root         string         // FS root path
    FSKind       string         // Filesystem type (local, git, bundle)
    FSName       string         // Filesystem name
    NodeKind     NodeKind       // New: root, file, dir, bundle, bundled
    ParentID     *string        // Hierarchical structure
    NestedAssets []*Asset       // Recursive nesting
    Thumbnail    *string        // Bubble-up logic
    Tags         []*Tag         // Many-to-many
    Properties   Properties     // Key-value map
}
```

**Key Changes:**
1. **String IDs** - MD5 hash-based instead of auto-increment
2. **Hierarchical structure** - Parent/child relationships
3. **Filesystem awareness** - FSKind, FSName, Root fields
4. **Node types** - Explicit NodeKind for different asset types
5. **Thumbnail bubbling** - AfterSave hook propagates thumbnails up hierarchy
6. **Properties** - Generic key-value storage instead of fixed fields

---

### 2.3 API Architecture

#### V1 API Pattern (Echo)
```go
// core/api/projects/endpoints.go
func GetProject(c echo.Context) error {
    id := c.Param("id")
    // Direct GORM query
    project := database.GetProjectByID(id)
    return c.JSON(200, project)
}

// Registered in stlib.go
e.GET("/api/projects/:id", api.GetProject)
```

#### V2 API Pattern (Chi + Repository)
```go
// v2/library/api/asset.go
func (h *LibraryHandler) GetAsset(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    // Repository pattern
    asset, err := h.repo.FindByID(r.Context(), id)
    // JSON encoding
    api.JSON(w, asset)
}

// Registered via router
r.Route("/api/lib", func(r chi.Router) {
    r.Get("/{id}", handler.GetAsset)
})
```

**Key Improvements:**
1. **Repository pattern** - Abstracted data access
2. **Context propagation** - Proper context usage
3. **Error handling** - Structured error responses
4. **Middleware support** - Chi's composable middleware

---

### 2.4 Configuration Architecture

#### V1 Config (config.toml)
```toml
port = 8000
libraryPath = "./testdata"
maxRenderWorkers = 10
fileBlacklist = [".DS_Store"]
```

#### V2 Config (v2/config/)
```toml
[server]
port = 8000

[[library.filesystems]]
kind = 'local'
name = 'Library'
path = '/library'

[[library.filesystems]]
kind = 'gitfs'
name = 'RatRig Parts'
[library.filesystems.config]
url = 'https://github.com/Rat-Rig/RatRig-PrintedParts'

[library.assetTypes]
# Extensible asset type definitions
```

**V2 Advantages:**
- Multiple filesystem support
- Structured sections
- Extensible configuration
- Type-safe parsing via Viper

---

## 3. Integration Strategy Overview

### 3.1 Strategic Approach

We will use a **phased, incremental migration** strategy with the following principles:

1. **Coexistence First** - Run v1 and v2 code side-by-side
2. **Feature Flagging** - Toggle between v1/v2 implementations
3. **Backward Compatibility** - Maintain v1 API contracts
4. **Gradual Transition** - Migrate features incrementally
5. **Risk Mitigation** - Maintain rollback capability at every phase

### 3.2 Migration Phases (6 Phases, 16-20 weeks)

```
Phase 1: Foundation (2-3 weeks)
├─ Merge branches
├─ Resolve conflicts
├─ Dual build system
└─ CI/CD updates

Phase 2: Backend Migration (4-5 weeks)
├─ Integrate /v2 modules
├─ Virtual filesystem layer
├─ Repository pattern adoption
└─ API version coexistence

Phase 3: Frontend Integration (3-4 weeks)
├─ Embed React build
├─ API client integration
├─ Component migration
└─ Routing setup

Phase 4: Feature Parity (4-5 weeks)
├─ Migrate integrations (OctoPrint, Klipper)
├─ WebSocket events
├─ Discovery pipeline
└─ Desktop app refinement

Phase 5: Testing & Validation (2-3 weeks)
├─ E2E testing
├─ Performance testing
├─ Security audit
└─ User acceptance testing

Phase 6: Deployment & Rollout (1-2 weeks)
├─ Canary deployment
├─ Gradual rollout
├─ Monitoring
└─ Documentation
```

---

## 4. Phase 1: Foundation & Preparation

### 4.1 Objectives
- Create a merged codebase with both architectures
- Establish build system supporting both versions
- Update CI/CD pipelines
- Document architectural decisions

### 4.2 Tasks

#### 4.2.1 Branch Merge Strategy

```bash
# Create integration branch
git checkout -b integration/v2-react-merge origin/main

# Merge v2-react with strategy
git merge --no-commit --no-ff origin/v2-react

# Expected conflicts:
# - go.mod / go.sum (dependency versions)
# - .gitignore (new ignores for frontend)
# - README.md (documentation differences)
# - main.go (entry point changes)
```

**Conflict Resolution Plan:**
- **go.mod**: Accept v2-react (newer dependencies)
- **.gitignore**: Merge both (union of ignores)
- **README.md**: Create new merged README with migration guide
- **main.go**: Keep v1 main.go, rename v2 entry points in /cmd/

#### 4.2.2 Directory Structure After Merge

```
/home/user/agent/
├── main.go                    # V1 entry point (legacy)
├── stlib.go                   # V1 orchestrator (legacy)
├── cmd/                       # V2 entry points
│   ├── command/main.go        # V2 web server
│   ├── desktop/main.go        # V2 desktop app
│   └── desktopgen/main.go     # Code generation
├── core/                      # V1 backend (keep for gradual migration)
│   ├── api/
│   ├── processing/
│   ├── data/
│   └── ... (all v1 modules)
├── v2/                        # V2 backend (new architecture)
│   ├── library/
│   ├── config/
│   ├── database/
│   └── ... (all v2 modules)
├── frontend/                  # React application
│   ├── src/
│   ├── dist/
│   ├── package.json
│   └── vite.config.ts
├── config.toml                # Configuration
├── .air.toml                  # Hot reload config
├── go.mod                     # Merged dependencies
└── buildAndRun.sh             # Updated build script
```

#### 4.2.3 Build System Updates

**Update `buildAndRun.sh`:**
```bash
#!/bin/bash
set -e

# Build mode: v1 | v2 | both
MODE=${1:-v2}

echo "Building MMP Agent in mode: $MODE"

case $MODE in
  v1)
    echo "Building V1 (legacy)..."
    go build -o ./bin/mmp-v1 .
    ;;
  v2)
    echo "Building V2 (modern)..."
    # Frontend build
    cd frontend && npm install && npm run build && cd ..
    # Backend build
    go build -o ./bin/mmp-v2 ./cmd/command/
    ;;
  both)
    echo "Building both versions..."
    go build -o ./bin/mmp-v1 .
    cd frontend && npm install && npm run build && cd ..
    go build -o ./bin/mmp-v2 ./cmd/command/
    ;;
  *)
    echo "Unknown mode: $MODE"
    exit 1
    ;;
esac

echo "Build complete!"
```

#### 4.2.4 Go Module Resolution

**Unified go.mod:**
```go
module github.com/eduardooliveira/stLib

go 1.22

require (
    // V1 dependencies
    github.com/labstack/echo/v4 v4.11.4

    // V2 dependencies
    github.com/go-chi/chi/v5 v5.0.10
    github.com/go-chi/cors v1.2.1

    // Shared dependencies
    gorm.io/gorm v1.25.5
    gorm.io/driver/sqlite v1.5.4
    github.com/spf13/viper v1.17.0
    github.com/gorilla/websocket v1.5.1

    // Upgrade conflicts (use latest)
    golang.org/x/crypto v0.35.0
    golang.org/x/net v0.38.0
)
```

**Resolution Strategy:**
1. Keep all v1 dependencies (for backward compatibility)
2. Add v2 dependencies (Chi, new modules)
3. Upgrade shared dependencies to latest compatible versions
4. Test both v1 and v2 builds after merge

#### 4.2.5 Configuration Migration

**Create config adapter:**
```go
// v2/config/adapter.go
package config

import (
    v1config "github.com/eduardooliveira/stLib/core/runtime"
)

// AdaptV1Config converts v1 config to v2 format
func AdaptV1Config(v1 *v1config.Config) *Config {
    return &Config{
        Server: ServerConfig{
            Port: v1.Port,
            Hostname: v1.ServerHostname,
        },
        Library: LibraryConfig{
            Filesystems: []FilesystemConfig{
                {
                    Kind: "local",
                    Name: "Library",
                    Path: v1.LibraryPath,
                },
            },
        },
    }
}
```

### 4.3 Deliverables
- ✅ Merged branch with both architectures
- ✅ Updated build scripts
- ✅ Resolved go.mod conflicts
- ✅ CI/CD pipeline updates
- ✅ Migration documentation

### 4.4 Success Criteria
- Both v1 and v2 build successfully
- All tests pass
- No breaking changes to existing deployments
- Documentation complete

---

## 5. Phase 2: Backend Migration

### 5.1 Objectives
- Integrate v2 backend modules alongside v1
- Implement feature flags for gradual rollout
- Maintain API compatibility
- Migrate data models incrementally

### 5.2 Tasks

#### 5.2.1 Virtual Filesystem Layer Integration

**Step 1: Add LibFS interface to v1 entities**

```go
// core/entities/filesystem.go (new file)
package entities

import (
    "github.com/eduardooliveira/stLib/v2/library/entities"
)

// FSAdapter wraps v2 LibFS for v1 usage
type FSAdapter struct {
    libFS entities.LibFS
}

func NewFSAdapter(libFS entities.LibFS) *FSAdapter {
    return &FSAdapter{libFS: libFS}
}

func (a *FSAdapter) ReadFile(path string) ([]byte, error) {
    return fs.ReadFile(a.libFS, path)
}
```

**Step 2: Update v1 processing to use virtual FS**

```go
// core/processing/discovery/assetFlatDiscovery.go
func DiscoverAssets(projectPath string) error {
    // OLD: Direct filesystem access
    // files, _ := os.ReadDir(projectPath)

    // NEW: Virtual filesystem
    fs := entities.GetFilesystem(projectPath)
    files, _ := fs.ReadDir(projectPath)
    // ... rest of logic unchanged
}
```

#### 5.2.2 Repository Pattern Adoption

**Create repository interfaces for v1 entities:**

```go
// core/data/repository/project_repository.go
package repository

import (
    "context"
    "github.com/eduardooliveira/stLib/core/entities"
    "gorm.io/gorm"
)

type ProjectRepository interface {
    FindByID(ctx context.Context, id uint) (*entities.Project, error)
    FindAll(ctx context.Context) ([]*entities.Project, error)
    Create(ctx context.Context, project *entities.Project) error
    Update(ctx context.Context, project *entities.Project) error
    Delete(ctx context.Context, id uint) error
}

type projectRepository struct {
    db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
    return &projectRepository{db: db}
}

func (r *projectRepository) FindByID(ctx context.Context, id uint) (*entities.Project, error) {
    var project entities.Project
    err := r.db.WithContext(ctx).First(&project, id).Error
    return &project, err
}
```

**Update v1 API handlers to use repositories:**

```go
// core/api/projects/endpoints.go
func (h *ProjectHandler) GetProject(c echo.Context) error {
    id := c.Param("id")

    // OLD: Direct database access
    // project := database.GetProjectByID(id)

    // NEW: Repository pattern
    project, err := h.repo.FindByID(c.Request().Context(), id)
    if err != nil {
        return c.JSON(404, map[string]string{"error": "Project not found"})
    }

    return c.JSON(200, project)
}
```

#### 5.2.3 API Versioning Strategy

**Implement dual API support:**

```go
// main.go (updated)
package main

import (
    "github.com/eduardooliveira/stLib/core"
    "github.com/eduardooliveira/stLib/v2"
)

func main() {
    // Feature flag
    enableV2 := os.Getenv("ENABLE_V2_API") == "true"

    if enableV2 {
        // Chi router with both APIs
        r := chi.NewRouter()

        // V1 API (compatibility)
        v1Handler := core.NewHandler()
        r.Mount("/api/v1", v1Handler)

        // V2 API (new)
        v2Handler := v2.NewLibraryHandler()
        r.Mount("/api/v2", v2Handler)

        // Legacy routes (redirect to v1)
        r.Mount("/api/projects", v1Handler.ProjectsRouter())

    } else {
        // Legacy v1 only
        core.Run()
    }
}
```

**API routing matrix:**

| Endpoint | V1 (Legacy) | V2 (New) | Notes |
|----------|-------------|----------|-------|
| `/api/projects` | ✅ Active | 🔄 Mapped to `/api/v2/lib` | Compatibility layer |
| `/api/projects/:id` | ✅ Active | 🔄 Mapped to `/api/v2/lib/:id` | Entity mapping |
| `/api/v1/projects` | ✅ Explicit v1 | - | For backward compat |
| `/api/v2/lib` | - | ✅ New API | V2 native |
| `/api/events` | ✅ WebSocket | ⏳ TODO | Needs migration |

#### 5.2.4 Data Migration Strategy

**Create migration utilities:**

```go
// v2/database/migration/v1_to_v2.go
package migration

import (
    v1entities "github.com/eduardooliveira/stLib/core/entities"
    v2entities "github.com/eduardooliveira/stLib/v2/library/entities"
)

// MigrateAsset converts v1 Asset to v2 Asset
func MigrateAsset(v1Asset *v1entities.Asset, fs v2entities.LibFS) *v2entities.Asset {
    return &v2entities.Asset{
        ID:          generateID(v1Asset.Path, fs),
        Label:       &v1Asset.Name,
        Path:        &v1Asset.Path,
        Root:        fs.GetRoot(),
        FSKind:      fs.Kind(),
        FSName:      fs.GetName(),
        // Map other fields...
    }
}

// MigrateDatabase performs full database migration
func MigrateDatabase(v1DB, v2DB *gorm.DB) error {
    // 1. Migrate Projects → Root Assets
    // 2. Migrate Assets → Nested Assets
    // 3. Migrate Tags (compatible)
    // 4. Migrate AssetTypes → Properties
    return nil
}
```

**Migration command:**

```bash
# Add to cmd/command/main.go
go run ./cmd/command/ --migrate-v1-data --v1-db=./data/v1.db --v2-db=./data/v2.db
```

#### 5.2.5 Database Schema Coexistence

**Option A: Dual databases (recommended for initial phase)**
```toml
# config.toml
[database]
v1_path = "./data/v1.db"  # Legacy SQLite
v2_path = "./data/v2.db"  # New SQLite

# Sync mode: write to both, read from v2
sync_mode = true
```

**Option B: Shared database with table prefixes**
```go
// v2/database/init.go
func Init() error {
    db.AutoMigrate(
        &entities.Asset{},      // v2_assets table
        &entities.Tag{},        // v2_tags table
    )

    // Keep v1 tables as-is for compatibility
    // - projects
    // - assets
    // - tags
}
```

### 5.3 Deliverables
- ✅ Virtual filesystem integration in v1 codebase
- ✅ Repository pattern implemented for v1 entities
- ✅ Dual API routing (v1 + v2)
- ✅ Data migration utilities
- ✅ Feature flag system

### 5.4 Success Criteria
- V1 API continues to work unchanged
- V2 API accessible via feature flag
- Data migration tools tested with sample data
- No performance degradation

---

## 6. Phase 3: Frontend Integration

### 6.1 Objectives
- Embed React frontend into application
- Integrate with backend APIs
- Migrate UI features from v1 (if any)
- Establish component library

### 6.2 Tasks

#### 6.2.1 Frontend Build Integration

**Update build pipeline:**

```yaml
# .github/workflows/build.yml
name: Build MMP Agent

on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      # Frontend build
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '20'

      - name: Build Frontend
        run: |
          cd frontend
          npm ci
          npm run build

      # Backend build
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.22'

      - name: Build Backend
        run: |
          go build -o ./bin/mmp ./cmd/command/

      - name: Run Tests
        run: |
          go test ./...
```

#### 6.2.2 Frontend Embedding

**Ensure Go embed directive:**

```go
// frontend/embeded.go
package frontend

import "embed"

//go:embed dist/*
var FS embed.FS
```

**Serve frontend in main:**

```go
// cmd/command/main.go (already implemented in v2-react)
assets, _ := fs.Sub(frontend.FS, "dist")
r.Handle("/assets/*", http.FileServerFS(assets))

// SPA fallback for React Router
r.NotFound(func(w http.ResponseWriter, r *http.Request) {
    http.ServeFileFS(w, r, frontend.FS, "/dist/index.html")
})
```

#### 6.2.3 API Client Integration

**Create API client:**

```typescript
// frontend/src/lib/api/client.ts
export class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string = '/api/v2') {
    this.baseUrl = baseUrl;
  }

  async getAsset(id: string): Promise<Asset> {
    const response = await fetch(`${this.baseUrl}/lib/${id}`);
    if (!response.ok) throw new Error('Failed to fetch asset');
    return response.json();
  }

  async listAssets(): Promise<Asset[]> {
    const response = await fetch(`${this.baseUrl}/lib`);
    return response.json();
  }
}

// Singleton instance
export const apiClient = new ApiClient();
```

**React Query integration:**

```typescript
// frontend/src/lib/fetchers/getAsset.ts
import { useQuery } from '@tanstack/react-query';
import { apiClient } from '../api/client';

export const useAsset = (id: string) => {
  return useQuery({
    queryKey: ['asset', id],
    queryFn: () => apiClient.getAsset(id),
  });
};

export const useAssets = () => {
  return useQuery({
    queryKey: ['assets'],
    queryFn: () => apiClient.listAssets(),
  });
};
```

#### 6.2.4 Component Migration Strategy

**Priority order for component migration:**

1. **Core Navigation** (Week 1)
   - NavBar
   - Sidebar
   - Routing setup

2. **Asset Management** (Week 2)
   - Asset list/grid view
   - Asset detail page
   - AssetCard component

3. **Forms & Editing** (Week 3)
   - Asset edit form
   - Asset creation form
   - Tag management

4. **3D Viewer** (Week 4)
   - Three.js integration
   - STL rendering
   - G-code visualization

**Component compatibility matrix:**

| Feature | V1 Status | V2 Component | Migration Priority |
|---------|-----------|--------------|-------------------|
| Project List | ❓ Unknown | `LibMain.tsx` | P0 (Critical) |
| Project Detail | ❓ Unknown | `Asset.tsx` | P0 (Critical) |
| Asset Grid | ❓ Unknown | `AssetCard.tsx` | P1 (High) |
| 3D Viewer | ❓ Unknown | `Viewer3D.tsx` | P1 (High) |
| Forms | ❓ Unknown | `AssetEditForm.tsx` | P2 (Medium) |
| Settings | ❓ Unknown | `LibSettings.tsx` | P3 (Low) |

#### 6.2.5 State Management

**Context setup:**

```typescript
// frontend/src/lib/contexts/AssetContext.tsx
import { createContext, useContext, ReactNode } from 'react';

interface AssetContextType {
  currentAsset: Asset | null;
  setCurrentAsset: (asset: Asset | null) => void;
  breadcrumbs: Asset[];
}

const AssetContext = createContext<AssetContextType | undefined>(undefined);

export const AssetProvider = ({ children }: { children: ReactNode }) => {
  const [currentAsset, setCurrentAsset] = useState<Asset | null>(null);

  // Calculate breadcrumbs from current asset
  const breadcrumbs = useMemo(() => {
    if (!currentAsset) return [];
    // Build path from root to current asset
    return buildBreadcrumbs(currentAsset);
  }, [currentAsset]);

  return (
    <AssetContext.Provider value={{ currentAsset, setCurrentAsset, breadcrumbs }}>
      {children}
    </AssetContext.Provider>
  );
};

export const useAssetContext = () => {
  const context = useContext(AssetContext);
  if (!context) throw new Error('useAssetContext must be used within AssetProvider');
  return context;
};
```

### 6.3 Deliverables
- ✅ Frontend build integrated into CI/CD
- ✅ React app embedded in Go binary
- ✅ API client with React Query
- ✅ Core components migrated
- ✅ Routing configured

### 6.4 Success Criteria
- Frontend accessible at `http://localhost:8000/`
- API calls successful to backend
- All core components render correctly
- No console errors

---

## 7. Phase 4: Feature Parity

### 7.1 Objectives
- Migrate v1-specific features to v2
- Implement WebSocket events
- Integrate hardware integrations
- Complete desktop application

### 7.2 Tasks

#### 7.2.1 WebSocket Events Migration

**V1 WebSocket implementation:**
```go
// core/events/events.go
var upgrader = websocket.Upgrader{}

func HandleWebSocket(c echo.Context) error {
    ws, _ := upgrader.Upgrade(c.Response(), c.Request(), nil)
    // Register client
    // Send events
}
```

**V2 WebSocket implementation:**

```go
// v2/events/manager.go
package events

import (
    "github.com/gorilla/websocket"
    "net/http"
)

type EventManager struct {
    clients   map[*websocket.Conn]bool
    broadcast chan Event
}

func NewEventManager() *EventManager {
    return &EventManager{
        clients:   make(map[*websocket.Conn]bool),
        broadcast: make(chan Event, 100),
    }
}

func (m *EventManager) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, _ := upgrader.Upgrade(w, r, nil)
    m.clients[conn] = true

    // Listen for events
    go m.listenForEvents(conn)
}

func (m *EventManager) BroadcastAssetUpdate(asset *entities.Asset) {
    m.broadcast <- Event{
        Type: "asset.updated",
        Data: asset,
    }
}
```

**Register WebSocket route:**

```go
// cmd/command/main.go
eventMgr := events.NewEventManager()
go eventMgr.Start()

r.Get("/api/events", eventMgr.HandleWebSocket)
```

**Frontend WebSocket client:**

```typescript
// frontend/src/lib/api/websocket.ts
export class WebSocketClient {
  private ws: WebSocket | null = null;
  private listeners: Map<string, Set<(data: any) => void>> = new Map();

  connect() {
    this.ws = new WebSocket('ws://localhost:8000/api/events');

    this.ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      const handlers = this.listeners.get(message.type);
      handlers?.forEach(handler => handler(message.data));
    };
  }

  on(eventType: string, handler: (data: any) => void) {
    if (!this.listeners.has(eventType)) {
      this.listeners.set(eventType, new Set());
    }
    this.listeners.get(eventType)!.add(handler);
  }
}

export const wsClient = new WebSocketClient();
```

#### 7.2.2 Printer Integration Migration

**Map v1 printer integrations:**

```go
// v2/integrations/printers/ (new)
package printers

import (
    v1printers "github.com/eduardooliveira/stLib/core/integrations/printers"
)

// Adapter wraps v1 printer implementations
type PrinterAdapter struct {
    v1Printer *v1printers.Printer
}

func (a *PrinterAdapter) GetState() (*State, error) {
    v1State := a.v1Printer.GetState()
    return convertState(v1State), nil
}

// Copy v1 implementations:
// - /core/integrations/octorpint/ → /v2/integrations/octoprint/
// - /core/integrations/klipper/ → /v2/integrations/klipper/
```

#### 7.2.3 Discovery Pipeline Integration

**V2 discovery using virtual FS:**

```go
// v2/library/discovery/scanner.go
func (s *Scanner) ScanFilesystem(fs entities.LibFS) error {
    // Walk virtual filesystem
    err := fs.Walk(".", func(path string, d fs.DirEntry, err error) error {
        if d.IsDir() {
            asset := entities.NewAsset(fs, path, true, parent)
        } else {
            asset := entities.NewAsset(fs, path, false, parent)
        }

        // Save to repository
        s.repo.Save(asset)

        // Enqueue for processing
        s.processor.Enqueue(asset)
    })
}
```

**Processing pipeline:**

```go
// v2/library/process/processor.go
func (p *Processor) Process(asset *entities.Asset) error {
    // 1. Enrichment phase
    for _, enricher := range p.enrichers {
        if enricher.CanProcess(asset) {
            enricher.Enrich(asset)
        }
    }

    // 2. Rendering phase
    for _, renderer := range p.renderers {
        if renderer.CanRender(asset) {
            renderer.Render(asset)
        }
    }

    // 3. Update database
    return p.repo.Update(asset)
}
```

#### 7.2.4 Desktop Application Refinement

**Wails configuration:**

```go
// cmd/desktop/main.go
func main() {
    // Start HTTP server in background
    go func() {
        http.ListenAndServe(":8000", handler)
    }()

    // Create Wails application
    app := NewApp()

    err := wails.Run(&options.App{
        Title:     "Maker Management Platform",
        Width:     1024,
        Height:    768,
        MinWidth:  800,
        MinHeight: 600,
        AssetServer: &assetserver.Options{
            Assets: frontend.FS,
        },
        BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
        OnStartup:        app.startup,
        Bind: []interface{}{
            app,
        },
        Mac: &mac.Options{
            TitleBar: mac.TitleBarHiddenInset(),
            Appearance: mac.NSAppearanceNameDarkAqua,
            About: &mac.AboutInfo{
                Title: "MMP",
                Message: "Maker Management Platform v2",
            },
        },
    })

    if err != nil {
        println("Error:", err.Error())
    }
}
```

### 7.3 Deliverables
- ✅ WebSocket events working
- ✅ Printer integrations migrated
- ✅ Discovery pipeline functional
- ✅ Desktop app polished

### 7.4 Success Criteria
- All v1 features available in v2
- Real-time updates working
- Desktop app launches successfully
- No feature regressions

---

## 8. Phase 5: Testing & Validation

### 8.1 Testing Strategy

#### 8.1.1 Unit Tests

**Backend tests:**
```go
// v2/library/entities/asset_test.go
func TestNewAsset(t *testing.T) {
    fs := libfs.NewLocalFS("/library", "Library")
    asset := entities.NewAsset(fs, "model.stl", false, nil)

    assert.NotEmpty(t, asset.ID)
    assert.Equal(t, "model", *asset.Label)
    assert.Equal(t, ".stl", *asset.Extension)
}

// v2/library/repo/asset_test.go
func TestAssetRepository_FindByID(t *testing.T) {
    db := setupTestDB()
    repo := repo.NewAssetRepository(db)

    asset := createTestAsset()
    repo.Save(asset)

    found, err := repo.FindByID(asset.ID)
    assert.NoError(t, err)
    assert.Equal(t, asset.ID, found.ID)
}
```

**Frontend tests:**
```typescript
// frontend/src/lib/components/assetCard/AssetCard.test.tsx
import { render, screen } from '@testing-library/react';
import { AssetCard } from './AssetCard';

describe('AssetCard', () => {
  it('renders asset label', () => {
    const asset = { id: '1', label: 'Test Model', kind: 'stl' };
    render(<AssetCard asset={asset} />);
    expect(screen.getByText('Test Model')).toBeInTheDocument();
  });
});
```

#### 8.1.2 Integration Tests

```go
// tests/integration/api_test.go
func TestAssetAPI_EndToEnd(t *testing.T) {
    // Setup test server
    server := setupTestServer()
    defer server.Close()

    // Create asset
    resp := postJSON(server.URL+"/api/v2/lib", map[string]interface{}{
        "label": "Test Asset",
        "path": "/test.stl",
    })
    assert.Equal(t, 201, resp.StatusCode)

    // Get asset
    var asset entities.Asset
    json.NewDecoder(resp.Body).Decode(&asset)

    resp = getJSON(server.URL+"/api/v2/lib/"+asset.ID)
    assert.Equal(t, 200, resp.StatusCode)
}
```

#### 8.1.3 Performance Tests

```go
// tests/performance/load_test.go
func BenchmarkAssetList(b *testing.B) {
    repo := setupRepo()

    // Create 10000 assets
    for i := 0; i < 10000; i++ {
        repo.Save(createTestAsset())
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        assets, _ := repo.FindAll(context.Background())
        assert.NotEmpty(b, assets)
    }
}
```

**Performance targets:**
- Asset list API: < 100ms for 10,000 assets
- Asset detail API: < 50ms
- Discovery scan: < 5 seconds for 1,000 files
- Frontend initial load: < 2 seconds

#### 8.1.4 E2E Tests (Playwright)

```typescript
// tests/e2e/asset-management.spec.ts
import { test, expect } from '@playwright/test';

test('asset workflow', async ({ page }) => {
  await page.goto('http://localhost:8000');

  // Navigate to library
  await page.click('text=Library');

  // Create new asset
  await page.click('button:has-text("New Asset")');
  await page.fill('input[name="label"]', 'E2E Test Asset');
  await page.click('button:has-text("Save")');

  // Verify asset appears
  await expect(page.locator('text=E2E Test Asset')).toBeVisible();
});
```

### 8.2 Deliverables
- ✅ Unit test coverage > 80%
- ✅ Integration tests for all APIs
- ✅ Performance benchmarks
- ✅ E2E test suite

### 8.3 Success Criteria
- All tests passing
- Performance targets met
- No critical bugs

---

## 9. Phase 6: Deployment & Rollout

### 9.1 Deployment Strategy

#### 9.1.1 Docker Images

**Multi-stage Dockerfile:**

```dockerfile
# Build frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Build backend
FROM golang:1.22-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist
RUN go build -o /mmp ./cmd/command/

# Runtime
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=backend-builder /mmp ./
EXPOSE 8000
CMD ["./mmp"]
```

**Build images:**
```bash
# V2 image
docker build -t ghcr.io/maker-management-platform/agent:v2 .

# V2-beta (for canary)
docker build -t ghcr.io/maker-management-platform/agent:v2-beta .
```

#### 9.1.2 Rollout Phases

**Phase 1: Internal Testing (Week 1)**
- Deploy to staging environment
- Internal team testing
- Bug fixes

**Phase 2: Beta Release (Week 2)**
- Release as `v2-beta` tag
- Announce to community for testing
- Collect feedback
- Monitor metrics

**Phase 3: Canary Deployment (Week 3)**
- Deploy to 10% of users
- Monitor error rates
- Gradually increase to 50%

**Phase 4: Full Rollout (Week 4)**
- Deploy to all users
- Update `latest` tag to v2
- Announce v2 as stable

#### 9.1.3 Migration Guide for Users

**Create `MIGRATION_GUIDE.md`:**

```markdown
# Migration Guide: V1 to V2

## Prerequisites
- Backup your data folder
- Note your current config.toml settings

## Migration Steps

### Docker Users
```bash
# Backup data
cp -r ./data ./data.backup

# Update docker-compose.yml
# Change image tag to v2
image: ghcr.io/maker-management-platform/agent:v2

# Restart container
docker-compose down
docker-compose up -d
```

### Standalone Binary Users
```bash
# Download v2 binary
wget https://github.com/.../mmp-v2-linux-amd64

# Run migration (optional - will auto-migrate)
./mmp-v2 --migrate-v1-data

# Start v2
./mmp-v2
```

## Configuration Changes
- `libraryPath` → `library.filesystems[0].path`
- New filesystem types supported (git, bundles)

## Breaking Changes
- API v1 still available at `/api/v1/*`
- New API at `/api/v2/*`
- Desktop app now available

## Rollback
```bash
# Restore data backup
rm -rf ./data
mv ./data.backup ./data

# Use v1 image
docker-compose down
# Change image to v1
docker-compose up -d
```
```

### 9.2 Monitoring & Observability

**Add metrics:**
```go
// v2/utils/metrics.go
package utils

import (
    "github.com/prometheus/client_golang/prometheus"
)

var (
    AssetCount = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "mmp_assets_total",
        Help: "Total number of assets",
    })

    APIRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "mmp_api_request_duration_seconds",
            Help: "API request duration",
        },
        []string{"endpoint", "method"},
    )
)

func init() {
    prometheus.MustRegister(AssetCount)
    prometheus.MustRegister(APIRequestDuration)
}
```

**Expose metrics endpoint:**
```go
// cmd/command/main.go
import "github.com/prometheus/client_golang/prometheus/promhttp"

r.Handle("/metrics", promhttp.Handler())
```

### 9.3 Deliverables
- ✅ Docker images published
- ✅ Migration guide
- ✅ Monitoring dashboard
- ✅ Rollback procedures

### 9.4 Success Criteria
- Successful deployment with zero downtime
- < 1% error rate
- User feedback positive

---

## 10. Risks & Mitigation

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| **Data loss during migration** | Medium | Critical | • Mandatory backups<br>• Dual-database mode<br>• Automated migration tests |
| **API breaking changes** | High | High | • Maintain v1 API compatibility<br>• Version routing<br>• Deprecation notices |
| **Performance degradation** | Medium | High | • Performance benchmarks<br>• Load testing<br>• Caching strategies |
| **Frontend bundle size** | Low | Medium | • Code splitting<br>• Lazy loading<br>• Tree shaking |
| **Missing v1 features** | High | High | • Feature matrix documentation<br>• Phased migration<br>• Feature flags |
| **User adoption resistance** | Medium | Medium | • Clear documentation<br>• Migration support<br>• Beta program |
| **Desktop app issues** | Medium | Low | • Platform-specific testing<br>• Fallback to web version |
| **Merge conflicts** | High | Low | • Careful branch management<br>• Regular syncing<br>• Automated conflict detection |

---

## 11. Success Metrics

### 11.1 Technical Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| **Test Coverage** | > 80% | Go test coverage + Jest coverage |
| **API Response Time** | < 100ms (p95) | Prometheus metrics |
| **Frontend Load Time** | < 2s (initial load) | Lighthouse score |
| **Build Time** | < 5 minutes | CI/CD pipeline duration |
| **Docker Image Size** | < 100MB | Docker inspect |
| **Zero Critical Bugs** | 0 | GitHub Issues |

### 11.2 User Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| **Migration Success Rate** | > 95% | Telemetry (opt-in) |
| **User Adoption** | > 50% in 1 month | Version reporting |
| **Support Tickets** | < 10/week | GitHub Issues |
| **User Satisfaction** | > 4.5/5 | Discord feedback |

---

## 12. Rollback Plan

### 12.1 Rollback Triggers

- Error rate > 5%
- Data corruption detected
- Critical security vulnerability
- Performance degradation > 50%
- User adoption < 10% after 2 weeks

### 12.2 Rollback Procedure

**Step 1: Immediate Actions**
```bash
# Stop v2 deployment
docker-compose down

# Revert to v1
docker-compose.yml: image: ghcr.io/maker-management-platform/agent:v1

# Start v1
docker-compose up -d
```

**Step 2: Data Recovery**
```bash
# Restore from backup
cp -r ./data.backup/* ./data/

# Verify data integrity
sqlite3 ./data/v1.db "PRAGMA integrity_check;"
```

**Step 3: Communication**
- Announce rollback to users via Discord
- Post-mortem analysis
- Document issues
- Plan v2.1 with fixes

### 12.3 Post-Rollback Analysis

- Root cause analysis
- Bug fixes in separate branch
- Additional testing
- Re-plan rollout

---

## Appendix A: File Mapping Reference

### Backend Files Migration Map

| V1 File | V2 File | Status | Notes |
|---------|---------|--------|-------|
| `main.go` | `cmd/command/main.go` | ✅ Migrate | Entry point |
| `stlib.go` | `v2/library/library.go` | ✅ Refactor | Orchestration |
| `core/entities/asset.go` | `v2/library/entities/asset.go` | ✅ Rewrite | New model |
| `core/entities/project.go` | `v2/library/entities/asset.go` | ✅ Merge | Project = Root Asset |
| `core/api/projects/endpoints.go` | `v2/library/api/asset.go` | ✅ Migrate | API handlers |
| `core/processing/discovery/` | `v2/library/discovery/scanner.go` | ✅ Adapt | Virtual FS |
| `core/processing/enrichment/` | `v2/library/process/enrichers/` | ✅ Copy | Compatible |
| `core/data/database/` | `v2/database/` | ✅ Migrate | GORM upgrade |
| `core/integrations/octorpint/` | `v2/integrations/octoprint/` | 🔄 Copy | Direct copy |
| `core/integrations/klipper/` | `v2/integrations/klipper/` | 🔄 Copy | Direct copy |
| `core/events/` | `v2/events/` | 🔄 Adapt | WebSocket refactor |

### Frontend Files Migration Map

| Component | V2 Location | Priority | Dependencies |
|-----------|-------------|----------|--------------|
| **Navigation** | `frontend/src/core/components/navbar/` | P0 | React Router |
| **Asset List** | `frontend/src/lib/pages/libMain/` | P0 | React Query |
| **Asset Detail** | `frontend/src/lib/components/asset/` | P0 | Three.js |
| **Asset Card** | `frontend/src/lib/components/assetCard/` | P1 | Mantine |
| **3D Viewer** | `frontend/src/lib/components/viewer3d/` | P1 | Three.js |
| **Forms** | `frontend/src/lib/components/assetEditForm/` | P2 | Mantine Form |
| **Settings** | `frontend/src/lib/pages/libSettings/` | P3 | - |

---

## Appendix B: Configuration Schema Evolution

### V1 Config Schema
```toml
port = 8000
serverHostname = "localhost"
libraryPath = "./testdata"
maxRenderWorkers = 10
fileBlacklist = [".DS_Store"]
renderColor1 = "#167DF0"
renderColor2 = "#FFFFFF"
thingiverseToken = ""
```

### V2 Config Schema
```toml
[server]
port = 8000
hostname = "localhost"

[library]
maxRenderWorkers = 10

[[library.filesystems]]
kind = "local"
name = "Library"
path = "./testdata"

[[library.filesystems]]
kind = "gitfs"
name = "RemoteLib"
[library.filesystems.config]
url = "https://github.com/user/repo"

[library.rendering]
color1 = "#167DF0"
color2 = "#FFFFFF"

[library.assetTypes]
# Extensible asset type definitions

[library.blacklist]
files = [".DS_Store"]

[integrations.thingiverse]
token = ""
```

### Migration Function
```go
func MigrateConfig(v1 string) string {
    // Read v1 config
    // Transform to v2 schema
    // Write v2 config
}
```

---

## Appendix C: Timeline Gantt Chart

```
Week 1-3:   Phase 1 (Foundation)
            ████████████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░

Week 4-8:   Phase 2 (Backend Migration)
            ░░░░░░░░░░░░████████████████████░░░░░░░░░░

Week 9-12:  Phase 3 (Frontend Integration)
            ░░░░░░░░░░░░░░░░░░░░░░░░████████████░░░░░░

Week 13-17: Phase 4 (Feature Parity)
            ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░████████░░

Week 18-20: Phase 5 (Testing)
            ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░████░░

Week 21-22: Phase 6 (Deployment)
            ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░██
```

**Total Duration:** 20-22 weeks (~5-5.5 months)

---

## Appendix D: Team Roles & Responsibilities

| Role | Responsibilities | Time Commitment |
|------|------------------|-----------------|
| **Backend Lead** | Go migration, API design, data migration | 100% (20-22 weeks) |
| **Frontend Lead** | React integration, component development | 100% (12-22 weeks) |
| **DevOps Engineer** | CI/CD, Docker, deployment | 50% (ongoing) |
| **QA Engineer** | Test strategy, E2E tests, validation | 100% (weeks 15-22) |
| **Product Manager** | Feature prioritization, user communication | 25% (ongoing) |
| **Technical Writer** | Documentation, migration guide | 50% (weeks 18-22) |

---

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2025-11-06 | Claude AI | Initial draft |

---

## Approval

_This document requires approval from:_

- [ ] Project Lead
- [ ] Backend Lead
- [ ] Frontend Lead
- [ ] DevOps Lead

---

**Next Steps:**
1. Review and approve this strategy document
2. Assign team members to roles
3. Create GitHub project with milestones
4. Begin Phase 1 execution
5. Schedule weekly sync meetings

**Questions?** Open an issue on GitHub or discuss in Discord.
