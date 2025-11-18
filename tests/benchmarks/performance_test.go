package benchmarks

import (
	"context"
	"os"
	"testing"

	v1entities "github.com/eduardooliveira/stLib/core/entities"
	"github.com/eduardooliveira/stLib/core/integration"
	v1discovery "github.com/eduardooliveira/stLib/core/processing/discovery"
	v2events "github.com/eduardooliveira/stLib/v2/events"
	v2discovery "github.com/eduardooliveira/stLib/v2/library/discovery"
	v2entities "github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/library/libfs"
	v2process "github.com/eduardooliveira/stLib/v2/library/process"
	v2repo "github.com/eduardooliveira/stLib/v2/library/repo"
)

// BenchmarkV1AssetDiscovery benchmarks v1 asset discovery
func BenchmarkV1AssetDiscovery(b *testing.B) {
	tmpDir := b.TempDir()

	// Create some test files
	for i := 0; i < 10; i++ {
		f, _ := os.Create(tmpDir + "/test" + string(rune(i)) + ".stl")
		f.Close()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = v1discovery.AssetFlatDiscovery(tmpDir)
	}
}

// BenchmarkV2AssetDiscovery benchmarks v2 asset discovery
func BenchmarkV2AssetDiscovery(b *testing.B) {
	tmpDir := b.TempDir()

	// Create some test files
	for i := 0; i < 10; i++ {
		f, _ := os.Create(tmpDir + "/test" + string(rune(i)) + ".stl")
		f.Close()
	}

	fs := libfs.NewLocalFS(tmpDir, "Benchmark")
	repo := v2repo.NewInMemoryAssetRepository()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scanner := v2discovery.NewScanner(fs, repo)
		_, _ = scanner.Scan(context.Background())
	}
}

// BenchmarkDiscoveryComparison compares v1 vs v2 discovery
func BenchmarkDiscoveryComparison(b *testing.B) {
	tmpDir := b.TempDir()

	// Create test files
	for i := 0; i < 100; i++ {
		f, _ := os.Create(tmpDir + "/test" + string(rune(i)) + ".stl")
		f.Close()
	}

	b.Run("V1Discovery", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = v1discovery.AssetFlatDiscovery(tmpDir)
		}
	})

	b.Run("V2Discovery", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fs := libfs.NewLocalFS(tmpDir, "Benchmark")
			repo := v2repo.NewInMemoryAssetRepository()
			scanner := v2discovery.NewScanner(fs, repo)
			_, _ = scanner.Scan(context.Background())
		}
	})
}

// BenchmarkV2Processing benchmarks v2 asset processing
func BenchmarkV2Processing(b *testing.B) {
	tmpDir := b.TempDir()
	fs := libfs.NewLocalFS(tmpDir, "Benchmark")

	// Create test asset
	testFile := tmpDir + "/test.stl"
	os.WriteFile(testFile, []byte("test data"), 0644)

	processor := v2process.NewProcessor()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		asset := v2entities.NewAsset(fs, "test.stl", false, nil)
		_ = processor.Process(asset)
	}
}

// BenchmarkProcessingAdapter benchmarks processing adapter
func BenchmarkProcessingAdapter(b *testing.B) {
	tmpDir := b.TempDir()
	testFile := tmpDir + "/test.stl"
	os.WriteFile(testFile, []byte("test data"), 0644)

	b.Run("V1Mode", func(b *testing.B) {
		os.Clearenv()
		flags := integration.LoadFeatureFlags()
		adapter := integration.NewProcessingAdapter(flags, nil)

		v1Asset := &v1entities.Asset{
			Name: "test.stl",
			Path: testFile,
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = adapter.ProcessAsset(v1Asset)
		}
	})

	b.Run("V2Mode", func(b *testing.B) {
		os.Setenv("ENABLE_V2_PROCESSING", "true")
		defer os.Clearenv()

		flags := integration.LoadFeatureFlags()
		adapter := integration.NewProcessingAdapter(flags, nil)

		fs := libfs.NewLocalFS(tmpDir, "Benchmark")
		v2Asset := v2entities.NewAsset(fs, "test.stl", false, nil)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = adapter.ProcessAsset(v2Asset)
		}
	})
}

// BenchmarkDiscoveryAdapter benchmarks discovery adapter
func BenchmarkDiscoveryAdapter(b *testing.B) {
	tmpDir := b.TempDir()

	// Create test files
	for i := 0; i < 10; i++ {
		f, _ := os.Create(tmpDir + "/test" + string(rune(i)) + ".stl")
		f.Close()
	}

	ctx := context.Background()

	b.Run("V1Mode", func(b *testing.B) {
		os.Clearenv()
		flags := integration.LoadFeatureFlags()
		adapter := integration.NewDiscoveryAdapter(flags, nil, nil)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = adapter.DiscoverAssets(ctx, tmpDir)
		}
	})

	b.Run("V2Mode", func(b *testing.B) {
		os.Setenv("ENABLE_V2_PROCESSING", "true")
		os.Setenv("ENABLE_V2_LIBFS", "true")
		defer os.Clearenv()

		flags := integration.LoadFeatureFlags()
		repo := v2repo.NewInMemoryAssetRepository()
		adapter := integration.NewDiscoveryAdapter(flags, nil, repo)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = adapter.DiscoverAssets(ctx, tmpDir)
		}
	})
}

// BenchmarkEventBroadcasting benchmarks event broadcasting
func BenchmarkEventBroadcasting(b *testing.B) {
	mgr := v2events.NewEventManager()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go mgr.Start(ctx)

	testData := map[string]interface{}{
		"id":   "test-123",
		"name": "test.stl",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mgr.Broadcast(v2events.EventAssetCreated, testData)
	}
}

// BenchmarkAssetCreation benchmarks asset creation
func BenchmarkAssetCreation(b *testing.B) {
	tmpDir := b.TempDir()
	fs := libfs.NewLocalFS(tmpDir, "Benchmark")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v2entities.NewAsset(fs, "test.stl", false, nil)
	}
}

// BenchmarkIDGeneration benchmarks MD5 ID generation
func BenchmarkIDGeneration(b *testing.B) {
	migration := integration.NewDataMigration(nil, nil)
	path := "/test/path/to/file.stl"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = migration.GenerateAssetID(path)
	}
}

// BenchmarkAssetMigration benchmarks v1 to v2 asset migration
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

// BenchmarkProjectMigration benchmarks v1 to v2 project migration
func BenchmarkProjectMigration(b *testing.B) {
	tmpDir := b.TempDir()
	fs := libfs.NewLocalFS(tmpDir, "Benchmark")
	migration := integration.NewDataMigration(nil, nil)

	v1Project := &v1entities.Project{
		Name: "Test Project",
		Path: tmpDir,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = migration.MigrateProject(v1Project, fs)
	}
}

// BenchmarkLibFSAdapter benchmarks LibFS adapter operations
func BenchmarkLibFSAdapter(b *testing.B) {
	tmpDir := b.TempDir()

	// Create test file
	testFile := tmpDir + "/test.txt"
	os.WriteFile(testFile, []byte("test content"), 0644)

	adapter := integration.NewLibFSAdapter(tmpDir)

	b.Run("ReadFile", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = adapter.ReadFile("test.txt")
		}
	})

	b.Run("ReadDir", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = adapter.ReadDir(".")
		}
	})
}

// BenchmarkFeatureFlagLoading benchmarks feature flag loading
func BenchmarkFeatureFlagLoading(b *testing.B) {
	os.Clearenv()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = integration.LoadFeatureFlags()
	}
}

// BenchmarkAPIRouter benchmarks API router creation
func BenchmarkAPIRouter(b *testing.B) {
	os.Clearenv()
	flags := integration.LoadFeatureFlags()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = integration.NewAPIRouter(flags)
	}
}

// BenchmarkMemoryUsage benchmarks memory usage
func BenchmarkMemoryUsage(b *testing.B) {
	tmpDir := b.TempDir()

	// Create many test files
	for i := 0; i < 1000; i++ {
		f, _ := os.Create(tmpDir + "/test" + string(rune(i)) + ".stl")
		f.WriteString("test data")
		f.Close()
	}

	b.Run("V1Discovery", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = v1discovery.AssetFlatDiscovery(tmpDir)
		}
	})

	b.Run("V2Discovery", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			fs := libfs.NewLocalFS(tmpDir, "Benchmark")
			repo := v2repo.NewInMemoryAssetRepository()
			scanner := v2discovery.NewScanner(fs, repo)
			_, _ = scanner.Scan(context.Background())
		}
	})
}

// BenchmarkConcurrentDiscovery benchmarks concurrent discovery
func BenchmarkConcurrentDiscovery(b *testing.B) {
	tmpDir := b.TempDir()

	// Create test files
	for i := 0; i < 50; i++ {
		f, _ := os.Create(tmpDir + "/test" + string(rune(i)) + ".stl")
		f.Close()
	}

	os.Setenv("ENABLE_V2_PROCESSING", "true")
	os.Setenv("ENABLE_V2_LIBFS", "true")
	defer os.Clearenv()

	flags := integration.LoadFeatureFlags()
	repo := v2repo.NewInMemoryAssetRepository()
	adapter := integration.NewDiscoveryAdapter(flags, nil, repo)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = adapter.DiscoverAssets(context.Background(), tmpDir)
		}
	})
}
