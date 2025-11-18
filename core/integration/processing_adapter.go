package integration

import (
	"context"
	"log/slog"

	v1entities "github.com/eduardooliveira/stLib/core/entities"
	v1processing "github.com/eduardooliveira/stLib/core/processing"
	v2entities "github.com/eduardooliveira/stLib/v2/library/entities"
	v2process "github.com/eduardooliveira/stLib/v2/library/process"
	v2events "github.com/eduardooliveira/stLib/v2/events"
)

// ProcessingAdapter bridges v1 and v2 processing pipelines
type ProcessingAdapter struct {
	flags    *FeatureFlags
	eventMgr *v2events.EventManager
}

// NewProcessingAdapter creates a new processing adapter
func NewProcessingAdapter(flags *FeatureFlags, eventMgr *v2events.EventManager) *ProcessingAdapter {
	return &ProcessingAdapter{
		flags:    flags,
		eventMgr: eventMgr,
	}
}

// ProcessAsset processes an asset using the appropriate pipeline
func (a *ProcessingAdapter) ProcessAsset(asset interface{}) error {
	if a.flags.EnableV2Processing {
		return a.processV2Asset(asset)
	}
	return a.processV1Asset(asset)
}

// processV1Asset processes using v1 pipeline
func (a *ProcessingAdapter) processV1Asset(asset interface{}) error {
	v1Asset, ok := asset.(*v1entities.Asset)
	if !ok {
		return nil // Not a v1 asset, skip
	}

	slog.Debug("Processing asset with v1 pipeline", "asset", v1Asset.Name)

	// Use v1 processing
	if err := v1processing.ProcessAsset(v1Asset); err != nil {
		slog.Error("V1 processing failed", "error", err)
		return err
	}

	return nil
}

// processV2Asset processes using v2 pipeline
func (a *ProcessingAdapter) processV2Asset(asset interface{}) error {
	v2Asset, ok := asset.(*v2entities.Asset)
	if !ok {
		return nil // Not a v2 asset, skip
	}

	slog.Debug("Processing asset with v2 pipeline", "asset", *v2Asset.Label)

	// Use v2 processing
	processor := v2process.NewProcessor()
	if err := processor.Process(v2Asset); err != nil {
		slog.Error("V2 processing failed", "error", err)
		return err
	}

	// Broadcast processing completion if event manager available
	if a.eventMgr != nil {
		a.eventMgr.Broadcast(v2events.EventProcessCompleted, map[string]interface{}{
			"assetId": v2Asset.ID,
			"label":   v2Asset.Label,
		})
	}

	return nil
}

// ProcessFolder processes a folder using the appropriate pipeline
func (a *ProcessingAdapter) ProcessFolder(path string) error {
	ctx := context.Background()

	if a.flags.EnableV2Processing {
		return a.processV2Folder(ctx, path)
	}
	return a.processV1Folder(path)
}

// processV1Folder processes a folder using v1 pipeline
func (a *ProcessingAdapter) processV1Folder(path string) error {
	slog.Info("Processing folder with v1 pipeline", "path", path)

	// Broadcast scan started
	if a.eventMgr != nil {
		a.eventMgr.Broadcast(v2events.EventScanStarted, map[string]interface{}{
			"path": path,
		})
	}

	// Use v1 processing
	if err := v1processing.ProcessFolder(path); err != nil {
		slog.Error("V1 folder processing failed", "error", err)
		return err
	}

	// Broadcast scan completed
	if a.eventMgr != nil {
		a.eventMgr.Broadcast(v2events.EventScanCompleted, map[string]interface{}{
			"path": path,
		})
	}

	return nil
}

// processV2Folder processes a folder using v2 pipeline
func (a *ProcessingAdapter) processV2Folder(ctx context.Context, path string) error {
	slog.Info("Processing folder with v2 pipeline", "path", path)

	// Broadcast scan started
	if a.eventMgr != nil {
		a.eventMgr.Broadcast(v2events.EventScanStarted, map[string]interface{}{
			"path": path,
		})
	}

	// Use v2 processing with LibFS adapter
	adapter := NewLibFSAdapter(path)
	libFS := adapter.GetLibFS()

	// Create v2 processor
	processor := v2process.NewProcessor()

	// Walk filesystem and process
	err := adapter.Walk(".", func(filePath string, info any, err error) error {
		if err != nil {
			return err
		}

		// Create asset from file
		asset := v2entities.NewAsset(libFS, filePath, false, nil)

		// Process asset
		if err := processor.Process(asset); err != nil {
			slog.Warn("Failed to process asset", "path", filePath, "error", err)
		}

		// Broadcast progress
		if a.eventMgr != nil {
			a.eventMgr.BroadcastAssetCreated(asset.ID, asset)
		}

		return nil
	})

	if err != nil {
		slog.Error("V2 folder processing failed", "error", err)
		return err
	}

	// Broadcast scan completed
	if a.eventMgr != nil {
		a.eventMgr.Broadcast(v2events.EventScanCompleted, map[string]interface{}{
			"path": path,
		})
	}

	return nil
}

// ConvertV1AssetToV2 converts a v1 asset to v2 format for processing
func (a *ProcessingAdapter) ConvertV1AssetToV2(v1Asset *v1entities.Asset, fs v2entities.LibFS) (*v2entities.Asset, error) {
	migration := NewDataMigration(nil, nil)
	return migration.MigrateAsset(v1Asset, fs, nil)
}

// ConvertV2AssetToV1 converts a v2 asset to v1 format (for backward compatibility)
func (a *ProcessingAdapter) ConvertV2AssetToV1(v2Asset *v2entities.Asset) (*v1entities.Asset, error) {
	v1Asset := &v1entities.Asset{
		Name: "",
		Path: "",
	}

	if v2Asset.Label != nil {
		v1Asset.Name = *v2Asset.Label
	}

	if v2Asset.Path != nil {
		v1Asset.Path = *v2Asset.Path
	}

	if v2Asset.Thumbnail != nil {
		v1Asset.Thumbnail = *v2Asset.Thumbnail
	}

	// Convert properties back to individual fields
	if renderedImage, ok := v2Asset.Properties["renderedImage"].(string); ok {
		v1Asset.RenderedImage = renderedImage
	}

	return v1Asset, nil
}
