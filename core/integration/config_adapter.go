package integration

import (
	v1config "github.com/eduardooliveira/stLib/core/runtime"
	v2config "github.com/eduardooliveira/stLib/v2/config"
)

// ConfigAdapter converts between v1 and v2 configurations
type ConfigAdapter struct{}

// NewConfigAdapter creates a new config adapter
func NewConfigAdapter() *ConfigAdapter {
	return &ConfigAdapter{}
}

// V1ToV2 converts v1 configuration to v2 format
func (a *ConfigAdapter) V1ToV2(v1 *v1config.Config) *v2config.Config {
	cfg := &v2config.Config{
		Server: v2config.ServerConfig{
			Port:     v1.Port,
			Hostname: v1.ServerHostname,
		},
		Library: v2config.LibraryConfig{
			MaxRenderWorkers: v1.MaxRenderWorkers,
			Filesystems: []v2config.FilesystemConfig{
				{
					Kind: "local",
					Name: "Library",
					Path: v1.LibraryPath,
				},
			},
		},
	}

	// Convert rendering config
	if v1.ModelRenderColor != "" || v1.ModelBackgroundColor != "" {
		cfg.Library.Rendering = v2config.RenderingConfig{
			Color1: v1.ModelRenderColor,
			Color2: v1.ModelBackgroundColor,
		}
	}

	// Convert blacklist
	if len(v1.FileBlacklist) > 0 {
		cfg.Library.Blacklist = v2config.BlacklistConfig{
			Files: v1.FileBlacklist,
		}
	}

	return cfg
}

// V2ToV1 converts v2 configuration to v1 format (for backward compatibility)
func (a *ConfigAdapter) V2ToV1(v2 *v2config.Config) *v1config.Config {
	cfg := &v1config.Config{
		Port:             v2.Server.Port,
		ServerHostname:   v2.Server.Hostname,
		MaxRenderWorkers: v2.Library.MaxRenderWorkers,
	}

	// Get first local filesystem as library path
	for _, fs := range v2.Library.Filesystems {
		if fs.Kind == "local" {
			cfg.LibraryPath = fs.Path
			break
		}
	}

	// Convert rendering config
	cfg.ModelRenderColor = v2.Library.Rendering.Color1
	cfg.ModelBackgroundColor = v2.Library.Rendering.Color2

	// Convert blacklist
	cfg.FileBlacklist = v2.Library.Blacklist.Files

	return cfg
}

// MergeConfigs merges v1 and v2 configurations, preferring v2 values
func (a *ConfigAdapter) MergeConfigs(v1 *v1config.Config, v2 *v2config.Config) *v2config.Config {
	// Start with v2 as base
	merged := *v2

	// If v2 has empty values, fill from v1
	if merged.Server.Port == 0 && v1.Port != 0 {
		merged.Server.Port = v1.Port
	}

	if merged.Server.Hostname == "" && v1.ServerHostname != "" {
		merged.Server.Hostname = v1.ServerHostname
	}

	if merged.Library.MaxRenderWorkers == 0 && v1.MaxRenderWorkers != 0 {
		merged.Library.MaxRenderWorkers = v1.MaxRenderWorkers
	}

	// Add v1 library path if v2 has no local filesystems
	hasLocalFS := false
	for _, fs := range merged.Library.Filesystems {
		if fs.Kind == "local" {
			hasLocalFS = true
			break
		}
	}

	if !hasLocalFS && v1.LibraryPath != "" {
		merged.Library.Filesystems = append(merged.Library.Filesystems, v2config.FilesystemConfig{
			Kind: "local",
			Name: "Library",
			Path: v1.LibraryPath,
		})
	}

	return &merged
}
