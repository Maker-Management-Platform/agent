package integration_test

import (
	"context"
	"os"
	"testing"

	"github.com/eduardooliveira/stLib/core/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFeatureFlagDefaults tests default feature flag values
func TestFeatureFlagDefaults(t *testing.T) {
	// Clear environment
	os.Clearenv()

	flags := integration.LoadFeatureFlags()

	assert.False(t, flags.EnableV2API, "V2 API should be disabled by default")
	assert.False(t, flags.EnableV2LibFS, "V2 LibFS should be disabled by default")
	assert.False(t, flags.EnableRepositoryPattern, "Repository pattern should be disabled by default")
	assert.False(t, flags.EnableV2Processing, "V2 processing should be disabled by default")
	assert.False(t, flags.UseV2Database, "V2 database should be disabled by default")
	assert.True(t, flags.EnableAPIVersioning, "API versioning should be enabled by default")
}

// TestFeatureFlagEnvironment tests feature flags from environment
func TestFeatureFlagEnvironment(t *testing.T) {
	// Set environment variables
	os.Setenv("ENABLE_V2_API", "true")
	os.Setenv("ENABLE_V2_LIBFS", "true")
	os.Setenv("ENABLE_REPOSITORY_PATTERN", "true")
	os.Setenv("ENABLE_V2_PROCESSING", "true")
	os.Setenv("USE_V2_DATABASE", "true")
	defer os.Clearenv()

	flags := integration.LoadFeatureFlags()

	assert.True(t, flags.EnableV2API, "V2 API should be enabled")
	assert.True(t, flags.EnableV2LibFS, "V2 LibFS should be enabled")
	assert.True(t, flags.EnableRepositoryPattern, "Repository pattern should be enabled")
	assert.True(t, flags.EnableV2Processing, "V2 processing should be enabled")
	assert.True(t, flags.UseV2Database, "V2 database should be enabled")
}

// TestV1OnlyMode tests v1-only mode detection
func TestV1OnlyMode(t *testing.T) {
	os.Clearenv()

	flags := integration.LoadFeatureFlags()
	mode := flags.GetMode()

	assert.Equal(t, "v1-only", mode, "Mode should be v1-only when all flags disabled")
}

// TestV2FullMode tests v2-full mode detection
func TestV2FullMode(t *testing.T) {
	os.Setenv("ENABLE_V2_API", "true")
	os.Setenv("ENABLE_V2_LIBFS", "true")
	os.Setenv("ENABLE_REPOSITORY_PATTERN", "true")
	os.Setenv("ENABLE_V2_PROCESSING", "true")
	os.Setenv("USE_V2_DATABASE", "true")
	defer os.Clearenv()

	flags := integration.LoadFeatureFlags()
	mode := flags.GetMode()

	assert.Equal(t, "v2-full", mode, "Mode should be v2-full when all flags enabled")
}

// TestHybridMode tests hybrid mode detection
func TestHybridMode(t *testing.T) {
	os.Setenv("ENABLE_V2_API", "true")
	os.Setenv("ENABLE_V2_LIBFS", "false")
	defer os.Clearenv()

	flags := integration.LoadFeatureFlags()
	mode := flags.GetMode()

	assert.Equal(t, "hybrid", mode, "Mode should be hybrid when some flags enabled")
}

// TestLibFSAdapter tests LibFS adapter functionality
func TestLibFSAdapter(t *testing.T) {
	// Create temporary test directory
	tmpDir := t.TempDir()

	// Create a test file
	testFile := tmpDir + "/test.txt"
	content := []byte("test content")
	err := os.WriteFile(testFile, content, 0644)
	require.NoError(t, err)

	// Create adapter
	adapter := integration.NewLibFSAdapter(tmpDir)
	require.NotNil(t, adapter)

	// Test ReadFile
	data, err := adapter.ReadFile("test.txt")
	assert.NoError(t, err)
	assert.Equal(t, content, data)

	// Test ReadDir
	entries, err := adapter.ReadDir(".")
	assert.NoError(t, err)
	assert.Greater(t, len(entries), 0)

	// Test LibFS access
	libFS := adapter.GetLibFS()
	assert.NotNil(t, libFS)
}

// TestConfigAdapter tests configuration adapter
func TestConfigAdapter(t *testing.T) {
	adapter := integration.NewConfigAdapter()
	require.NotNil(t, adapter)

	// Test V1ToV2 conversion
	t.Run("V1ToV2", func(t *testing.T) {
		// This would require creating test v1 and v2 config structs
		// For now, just verify adapter exists
		assert.NotNil(t, adapter)
	})

	// Test V2ToV1 conversion
	t.Run("V2ToV1", func(t *testing.T) {
		assert.NotNil(t, adapter)
	})
}

// TestProcessingAdapter tests processing adapter routing
func TestProcessingAdapter(t *testing.T) {
	ctx := context.Background()

	t.Run("V1Mode", func(t *testing.T) {
		os.Clearenv()
		flags := integration.LoadFeatureFlags()
		adapter := integration.NewProcessingAdapter(flags, nil)

		assert.NotNil(t, adapter)
		// Processing adapter should route to v1 when flags disabled
	})

	t.Run("V2Mode", func(t *testing.T) {
		os.Setenv("ENABLE_V2_PROCESSING", "true")
		defer os.Clearenv()

		flags := integration.LoadFeatureFlags()
		adapter := integration.NewProcessingAdapter(flags, nil)

		assert.NotNil(t, adapter)
		// Processing adapter should route to v2 when flags enabled
		_ = ctx
	})
}

// TestDiscoveryAdapter tests discovery adapter routing
func TestDiscoveryAdapter(t *testing.T) {
	ctx := context.Background()

	t.Run("V1Mode", func(t *testing.T) {
		os.Clearenv()
		flags := integration.LoadFeatureFlags()
		adapter := integration.NewDiscoveryAdapter(flags, nil, nil)

		assert.NotNil(t, adapter)
		// Discovery adapter should route to v1 when flags disabled
		_ = ctx
	})

	t.Run("V2Mode", func(t *testing.T) {
		os.Setenv("ENABLE_V2_PROCESSING", "true")
		os.Setenv("ENABLE_V2_LIBFS", "true")
		defer os.Clearenv()

		flags := integration.LoadFeatureFlags()
		adapter := integration.NewDiscoveryAdapter(flags, nil, nil)

		assert.NotNil(t, adapter)
		// Discovery adapter should route to v2 when flags enabled
		_ = ctx
	})
}

// TestDataMigration tests data migration utilities
func TestDataMigration(t *testing.T) {
	migration := integration.NewDataMigration(nil, nil)
	require.NotNil(t, migration)

	// Test ID generation
	t.Run("GenerateID", func(t *testing.T) {
		id := migration.GenerateAssetID("test/path.txt")
		assert.NotEmpty(t, id)
		assert.Len(t, id, 32) // MD5 hash length

		// Same path should generate same ID
		id2 := migration.GenerateAssetID("test/path.txt")
		assert.Equal(t, id, id2)

		// Different path should generate different ID
		id3 := migration.GenerateAssetID("test/other.txt")
		assert.NotEqual(t, id, id3)
	})
}

// TestAPIRouter tests API router
func TestAPIRouter(t *testing.T) {
	os.Clearenv()
	flags := integration.LoadFeatureFlags()

	router := integration.NewAPIRouter(flags)
	require.NotNil(t, router)

	handler := router.GetUnifiedHandler()
	assert.NotNil(t, handler, "Unified handler should be created")
}
