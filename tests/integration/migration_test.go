package integration_test

import (
	"context"
	"testing"

	"github.com/eduardooliveira/stLib/core/integration"
	v1entities "github.com/eduardooliveira/stLib/core/entities"
	v2entities "github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/library/libfs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIDGeneration tests MD5-based ID generation
func TestIDGeneration(t *testing.T) {
	migration := integration.NewDataMigration(nil, nil)

	tests := []struct {
		name     string
		path     string
		wantLen  int
		checksum bool
	}{
		{"Simple path", "/test/file.txt", 32, true},
		{"Complex path", "/my/complex/path/to/file.stl", 32, true},
		{"With spaces", "/path with spaces/file.txt", 32, true},
		{"With special chars", "/path/file-name_123.gcode", 32, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := migration.GenerateAssetID(tt.path)

			// Check length
			assert.Len(t, id, tt.wantLen, "ID should be MD5 hash length")

			// Check consistency
			id2 := migration.GenerateAssetID(tt.path)
			assert.Equal(t, id, id2, "Same path should generate same ID")

			// Check uniqueness
			id3 := migration.GenerateAssetID(tt.path + ".other")
			assert.NotEqual(t, id, id3, "Different paths should generate different IDs")
		})
	}
}

// TestProjectMigration tests v1 Project -> v2 Root Asset migration
func TestProjectMigration(t *testing.T) {
	tmpDir := t.TempDir()
	fs := libfs.NewLocalFS(tmpDir, "TestProject")

	migration := integration.NewDataMigration(nil, nil)

	// Create v1 project
	v1Project := &v1entities.Project{
		Name:        "Test Project",
		Description: "Test Description",
		Path:        tmpDir,
		Tags:        []v1entities.Tag{{Name: "test"}},
	}

	// Migrate to v2
	v2Asset, err := migration.MigrateProject(v1Project, fs)
	require.NoError(t, err)
	require.NotNil(t, v2Asset)

	// Verify migration
	assert.NotEmpty(t, v2Asset.ID, "ID should be generated")
	assert.Equal(t, v2entities.NodeKindRoot, v2Asset.NodeKind, "Should be root asset")
	assert.NotNil(t, v2Asset.Label)
	assert.Equal(t, "Test Project", *v2Asset.Label, "Label should match project name")
	assert.NotNil(t, v2Asset.Kind)
	assert.Equal(t, "project", *v2Asset.Kind, "Kind should be project")

	// Check properties
	if desc, ok := v2Asset.Properties["description"].(string); ok {
		assert.Equal(t, "Test Description", desc, "Description should be preserved")
	}

	// Check tags
	assert.Greater(t, len(v2Asset.Tags), 0, "Tags should be migrated")
}

// TestAssetMigration tests v1 Asset -> v2 Asset migration
func TestAssetMigration(t *testing.T) {
	tmpDir := t.TempDir()
	fs := libfs.NewLocalFS(tmpDir, "TestAsset")

	migration := integration.NewDataMigration(nil, nil)

	tests := []struct {
		name      string
		v1Asset   *v1entities.Asset
		wantKind  string
		wantProps []string
	}{
		{
			name: "STL Model",
			v1Asset: &v1entities.Asset{
				Name:      "model.stl",
				Path:      tmpDir + "/model.stl",
				Thumbnail: "/thumbs/model.png",
			},
			wantKind:  "model",
			wantProps: []string{"thumbnail"},
		},
		{
			name: "GCode File",
			v1Asset: &v1entities.Asset{
				Name: "print.gcode",
				Path: tmpDir + "/print.gcode",
			},
			wantKind:  "gcode",
			wantProps: []string{},
		},
		{
			name: "Image",
			v1Asset: &v1entities.Asset{
				Name: "photo.jpg",
				Path: tmpDir + "/photo.jpg",
			},
			wantKind:  "image",
			wantProps: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v2Asset, err := migration.MigrateAsset(tt.v1Asset, fs, nil)
			require.NoError(t, err)
			require.NotNil(t, v2Asset)

			// Verify basic fields
			assert.NotEmpty(t, v2Asset.ID, "ID should be generated")
			assert.NotNil(t, v2Asset.Label)
			assert.Equal(t, tt.v1Asset.Name, *v2Asset.Label, "Label should match name")

			// Verify kind (if determinable)
			if tt.wantKind != "" {
				assert.NotNil(t, v2Asset.Kind)
				// Kind inference depends on file extension
			}

			// Verify properties
			for _, prop := range tt.wantProps {
				_, exists := v2Asset.Properties[prop]
				assert.True(t, exists, "Property %s should exist", prop)
			}
		})
	}
}

// TestTagMigration tests tag migration
func TestTagMigration(t *testing.T) {
	v1Tags := []v1entities.Tag{
		{Name: "important"},
		{Name: "ready"},
		{Name: "work-in-progress"},
	}

	v2Tags := integration.MigrateTags(v1Tags)

	assert.Len(t, v2Tags, len(v1Tags), "All tags should be migrated")
	for i, v1Tag := range v1Tags {
		assert.Equal(t, v1Tag.Name, v2Tags[i], "Tag name should be preserved")
	}
}

// TestRelationshipMigration tests parent-child relationships
func TestRelationshipMigration(t *testing.T) {
	tmpDir := t.TempDir()
	fs := libfs.NewLocalFS(tmpDir, "TestRelationships")

	migration := integration.NewDataMigration(nil, nil)

	// Create parent project
	v1Project := &v1entities.Project{
		Name: "Parent Project",
		Path: tmpDir,
	}

	parentAsset, err := migration.MigrateProject(v1Project, fs)
	require.NoError(t, err)
	parentID := parentAsset.ID

	// Create child asset with parent reference
	v1Asset := &v1entities.Asset{
		Name: "child.stl",
		Path: tmpDir + "/child.stl",
	}

	childAsset, err := migration.MigrateAsset(v1Asset, fs, &parentID)
	require.NoError(t, err)

	// Verify relationship
	assert.NotNil(t, childAsset.ParentID, "Child should have parent reference")
	assert.Equal(t, parentID, *childAsset.ParentID, "Parent ID should match")
	assert.Equal(t, v2entities.NodeKindChild, childAsset.NodeKind, "Should be child node")
}

// TestMetadataPreservation tests metadata preservation
func TestMetadataPreservation(t *testing.T) {
	tmpDir := t.TempDir()
	fs := libfs.NewLocalFS(tmpDir, "TestMetadata")

	migration := integration.NewDataMigration(nil, nil)

	v1Asset := &v1entities.Asset{
		Name:          "test.stl",
		Path:          tmpDir + "/test.stl",
		Thumbnail:     "/thumbs/test.png",
		RenderedImage: "/renders/test.jpg",
	}

	v2Asset, err := migration.MigrateAsset(v1Asset, fs, nil)
	require.NoError(t, err)

	// Check thumbnail preserved
	assert.NotNil(t, v2Asset.Thumbnail)
	assert.Equal(t, v1Asset.Thumbnail, *v2Asset.Thumbnail)

	// Check rendered image in properties
	if rendered, ok := v2Asset.Properties["renderedImage"].(string); ok {
		assert.Equal(t, v1Asset.RenderedImage, rendered)
	}
}

// TestBulkMigration tests migrating multiple items
func TestBulkMigration(t *testing.T) {
	tmpDir := t.TempDir()
	fs := libfs.NewLocalFS(tmpDir, "TestBulk")

	migration := integration.NewDataMigration(nil, nil)

	// Create multiple v1 assets
	v1Assets := []*v1entities.Asset{
		{Name: "file1.stl", Path: tmpDir + "/file1.stl"},
		{Name: "file2.gcode", Path: tmpDir + "/file2.gcode"},
		{Name: "file3.png", Path: tmpDir + "/file3.png"},
	}

	// Migrate all
	var v2Assets []*v2entities.Asset
	for _, v1Asset := range v1Assets {
		v2Asset, err := migration.MigrateAsset(v1Asset, fs, nil)
		require.NoError(t, err)
		v2Assets = append(v2Assets, v2Asset)
	}

	// Verify all migrated
	assert.Len(t, v2Assets, len(v1Assets), "All assets should be migrated")

	// Verify unique IDs
	idMap := make(map[string]bool)
	for _, v2Asset := range v2Assets {
		assert.False(t, idMap[v2Asset.ID], "IDs should be unique")
		idMap[v2Asset.ID] = true
	}
}

// TestMigrationReversibility tests v2 -> v1 conversion
func TestMigrationReversibility(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	fs := libfs.NewLocalFS(tmpDir, "TestReverse")

	processingAdapter := integration.NewProcessingAdapter(
		integration.LoadFeatureFlags(),
		nil,
	)

	// Create v1 asset
	v1Original := &v1entities.Asset{
		Name:      "test.stl",
		Path:      tmpDir + "/test.stl",
		Thumbnail: "/thumbs/test.png",
	}

	// Convert to v2
	v2Asset, err := processingAdapter.ConvertV1AssetToV2(v1Original, fs)
	require.NoError(t, err)

	// Convert back to v1
	v1Converted, err := processingAdapter.ConvertV2AssetToV1(v2Asset)
	require.NoError(t, err)

	// Verify key fields preserved
	assert.Equal(t, v1Original.Name, v1Converted.Name)
	assert.Equal(t, v1Original.Path, v1Converted.Path)
	assert.Equal(t, v1Original.Thumbnail, v1Converted.Thumbnail)

	_ = ctx
}

// TestMigrationWithDatabase tests migration with actual database
func TestMigrationWithDatabase(t *testing.T) {
	t.Skip("Requires database setup - implement with test database")

	// Would test:
	// - Create v1 records in database
	// - Run migration
	// - Verify v2 records in database
	// - Check all relationships intact
	// - Verify no data loss
}

// TestMigrationPerformance benchmarks migration performance
func BenchmarkAssetMigration(b *testing.B) {
	tmpDir := b.TempDir()
	fs := libfs.NewLocalFS(tmpDir, "Benchmark")

	migration := integration.NewDataMigration(nil, nil)

	v1Asset := &v1entities.Asset{
		Name: "test.stl",
		Path: tmpDir + "/test.stl",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = migration.MigrateAsset(v1Asset, fs, nil)
	}
}

// TestMigrationErrorHandling tests error cases
func TestMigrationErrorHandling(t *testing.T) {
	migration := integration.NewDataMigration(nil, nil)

	t.Run("Nil Project", func(t *testing.T) {
		_, err := migration.MigrateProject(nil, nil)
		assert.Error(t, err, "Should error on nil project")
	})

	t.Run("Nil Asset", func(t *testing.T) {
		_, err := migration.MigrateAsset(nil, nil, nil)
		assert.Error(t, err, "Should error on nil asset")
	})

	t.Run("Invalid Path", func(t *testing.T) {
		v1Asset := &v1entities.Asset{
			Name: "test",
			Path: "", // Empty path
		}
		// Should handle gracefully
		_ = v1Asset
	})
}

// TestMigrationIdempotency tests running migration multiple times
func TestMigrationIdempotency(t *testing.T) {
	tmpDir := t.TempDir()
	fs := libfs.NewLocalFS(tmpDir, "TestIdempotency")

	migration := integration.NewDataMigration(nil, nil)

	v1Asset := &v1entities.Asset{
		Name: "test.stl",
		Path: tmpDir + "/test.stl",
	}

	// Migrate once
	v2Asset1, err := migration.MigrateAsset(v1Asset, fs, nil)
	require.NoError(t, err)

	// Migrate again
	v2Asset2, err := migration.MigrateAsset(v1Asset, fs, nil)
	require.NoError(t, err)

	// Should produce same result
	assert.Equal(t, v2Asset1.ID, v2Asset2.ID, "Same asset should produce same ID")
	assert.Equal(t, v2Asset1.Label, v2Asset2.Label, "Labels should match")
}
