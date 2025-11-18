package integration

import (
	"context"
	"log/slog"

	v1discovery "github.com/eduardooliveira/stLib/core/processing/discovery"
	v2discovery "github.com/eduardooliveira/stLib/v2/library/discovery"
	v2entities "github.com/eduardooliveira/stLib/v2/library/entities"
	v2repo "github.com/eduardooliveira/stLib/v2/library/repo"
	v2events "github.com/eduardooliveira/stLib/v2/events"
	"github.com/eduardooliveira/stLib/v2/library/libfs"
)

// DiscoveryAdapter bridges v1 and v2 discovery systems
type DiscoveryAdapter struct {
	flags    *FeatureFlags
	eventMgr *v2events.EventManager
	repo     *v2repo.AssetRepository
}

// NewDiscoveryAdapter creates a new discovery adapter
func NewDiscoveryAdapter(flags *FeatureFlags, eventMgr *v2events.EventManager, repo *v2repo.AssetRepository) *DiscoveryAdapter {
	return &DiscoveryAdapter{
		flags:    flags,
		eventMgr: eventMgr,
		repo:     repo,
	}
}

// DiscoverAssets discovers assets in a path using the appropriate system
func (a *DiscoveryAdapter) DiscoverAssets(ctx context.Context, path string) error {
	if a.flags.EnableV2Processing && a.flags.EnableV2LibFS {
		return a.discoverV2(ctx, path)
	}
	return a.discoverV1(path)
}

// discoverV1 uses v1 discovery system
func (a *DiscoveryAdapter) discoverV1(path string) error {
	slog.Info("Running v1 asset discovery", "path", path)

	// Broadcast discovery started
	if a.eventMgr != nil {
		a.eventMgr.Broadcast(v2events.EventScanStarted, map[string]interface{}{
			"path": path,
			"mode": "v1",
		})
	}

	// Use v1 flat discovery
	assets, err := v1discovery.AssetFlatDiscovery(path)
	if err != nil {
		slog.Error("V1 discovery failed", "error", err)
		return err
	}

	slog.Info("V1 discovery completed", "count", len(assets))

	// Broadcast discovery completed
	if a.eventMgr != nil {
		a.eventMgr.Broadcast(v2events.EventScanCompleted, map[string]interface{}{
			"path":  path,
			"count": len(assets),
			"mode":  "v1",
		})
	}

	return nil
}

// discoverV2 uses v2 discovery system
func (a *DiscoveryAdapter) discoverV2(ctx context.Context, path string) error {
	slog.Info("Running v2 asset discovery", "path", path)

	// Broadcast discovery started
	if a.eventMgr != nil {
		a.eventMgr.Broadcast(v2events.EventScanStarted, map[string]interface{}{
			"path": path,
			"mode": "v2",
		})
	}

	// Create LibFS for the path
	fs := libfs.NewLocalFS(path, "Discovery")

	// Create v2 scanner
	scanner := v2discovery.NewScanner(fs, a.repo)

	// Run discovery
	assets, err := scanner.Scan(ctx)
	if err != nil {
		slog.Error("V2 discovery failed", "error", err)
		return err
	}

	slog.Info("V2 discovery completed", "count", len(assets))

	// Broadcast each discovered asset
	if a.eventMgr != nil {
		for _, asset := range assets {
			a.eventMgr.BroadcastAssetCreated(asset.ID, asset)
		}

		a.eventMgr.Broadcast(v2events.EventScanCompleted, map[string]interface{}{
			"path":  path,
			"count": len(assets),
			"mode":  "v2",
		})
	}

	return nil
}

// DiscoverProjects discovers projects using the appropriate system
func (a *DiscoveryAdapter) DiscoverProjects(ctx context.Context, path string) error {
	if a.flags.EnableV2Processing && a.flags.EnableV2LibFS {
		// In v2, projects are root assets, use regular discovery
		return a.discoverV2(ctx, path)
	}

	// Use v1 deep discovery
	slog.Info("Running v1 project discovery", "path", path)

	// Broadcast discovery started
	if a.eventMgr != nil {
		a.eventMgr.Broadcast(v2events.EventScanStarted, map[string]interface{}{
			"path": path,
			"type": "projects",
			"mode": "v1",
		})
	}

	projects, err := v1discovery.ProjectDeepDiscovery(path)
	if err != nil {
		slog.Error("V1 project discovery failed", "error", err)
		return err
	}

	slog.Info("V1 project discovery completed", "count", len(projects))

	// Broadcast discovery completed
	if a.eventMgr != nil {
		a.eventMgr.Broadcast(v2events.EventScanCompleted, map[string]interface{}{
			"path":  path,
			"count": len(projects),
			"type":  "projects",
			"mode":  "v1",
		})
	}

	return nil
}

// MonitorChanges monitors a path for changes
func (a *DiscoveryAdapter) MonitorChanges(ctx context.Context, path string) error {
	// This would set up file system watchers
	// For now, it's a placeholder
	slog.Info("Setting up change monitoring", "path", path)
	return nil
}

// ScanMultiplePaths scans multiple paths
func (a *DiscoveryAdapter) ScanMultiplePaths(ctx context.Context, paths []string) error {
	slog.Info("Scanning multiple paths", "count", len(paths))

	for _, path := range paths {
		if err := a.DiscoverAssets(ctx, path); err != nil {
			slog.Error("Failed to scan path", "path", path, "error", err)
			// Continue with other paths
		}
	}

	return nil
}

// GetDiscoveryStats returns statistics about discovered assets
func (a *DiscoveryAdapter) GetDiscoveryStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	if a.repo != nil {
		// Get count from v2 repository
		assets, err := a.repo.FindAll(ctx)
		if err == nil {
			stats["totalAssets"] = len(assets)

			// Count by kind
			kindCounts := make(map[string]int)
			for _, asset := range assets {
				if asset.Kind != nil {
					kindCounts[*asset.Kind]++
				}
			}
			stats["byKind"] = kindCounts

			// Count by node kind
			nodeKindCounts := make(map[v2entities.NodeKind]int)
			for _, asset := range assets {
				nodeKindCounts[asset.NodeKind]++
			}
			stats["byNodeKind"] = nodeKindCounts
		}
	}

	return stats, nil
}
