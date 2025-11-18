package integration

import (
	"os"
	"strconv"
)

// FeatureFlags controls which features are enabled during v1→v2 migration
type FeatureFlags struct {
	// EnableV2API enables the new v2 API endpoints alongside v1
	EnableV2API bool

	// EnableV2LibFS enables the virtual filesystem layer
	EnableV2LibFS bool

	// EnableRepositoryPattern enables repository pattern for data access
	EnableRepositoryPattern bool

	// EnableV2Processing enables v2 processing pipeline
	EnableV2Processing bool

	// UseV2Database switches to v2 database schema
	UseV2Database bool

	// EnableAPIVersioning enables /api/v1 and /api/v2 routing
	EnableAPIVersioning bool
}

var (
	// Flags is the global feature flag instance
	Flags = LoadFeatureFlags()
)

// LoadFeatureFlags loads feature flags from environment variables
func LoadFeatureFlags() *FeatureFlags {
	return &FeatureFlags{
		EnableV2API:             getBoolEnv("ENABLE_V2_API", false),
		EnableV2LibFS:           getBoolEnv("ENABLE_V2_LIBFS", false),
		EnableRepositoryPattern: getBoolEnv("ENABLE_REPOSITORY_PATTERN", false),
		EnableV2Processing:      getBoolEnv("ENABLE_V2_PROCESSING", false),
		UseV2Database:           getBoolEnv("USE_V2_DATABASE", false),
		EnableAPIVersioning:     getBoolEnv("ENABLE_API_VERSIONING", true),
	}
}

// getBoolEnv gets a boolean environment variable with a default value
func getBoolEnv(key string, defaultValue bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}

	result, err := strconv.ParseBool(val)
	if err != nil {
		return defaultValue
	}

	return result
}

// IsV2Enabled returns true if any v2 feature is enabled
func (f *FeatureFlags) IsV2Enabled() bool {
	return f.EnableV2API || f.EnableV2LibFS || f.EnableRepositoryPattern ||
		f.EnableV2Processing || f.UseV2Database
}

// GetMode returns the current operating mode
func (f *FeatureFlags) GetMode() string {
	if !f.IsV2Enabled() {
		return "v1-only"
	}

	allV2 := f.EnableV2API && f.EnableV2LibFS && f.EnableRepositoryPattern &&
		f.EnableV2Processing && f.UseV2Database

	if allV2 {
		return "v2-full"
	}

	return "hybrid"
}
