package processing

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/eduardooliveira/stLib/core/entities"
	"github.com/eduardooliveira/stLib/core/processing/discovery"
	"github.com/eduardooliveira/stLib/core/processing/initialization"
	"github.com/eduardooliveira/stLib/core/runtime"
	"golang.org/x/sync/errgroup"
)

type ProcessableAsset struct {
	Name    string
	Label   string
	Project *entities.Project
	Asset   *entities.ProjectAsset
	Origin  string
}

func (p *ProcessableAsset) GetProject() *entities.Project {
	return p.Project
}
func (p *ProcessableAsset) GetAsset() *entities.ProjectAsset {
	return p.Asset
}

func ProcessFolder(ctx context.Context, root string) error {
	tempPath := filepath.Clean(filepath.Join(runtime.GetDataPath(), "assets")) //TODO: move this elsewhere
	if _, err := os.Stat(tempPath); os.IsNotExist(err) {
		err := os.MkdirAll(tempPath, os.ModePerm)
		if err != nil {
			log.Panic(err)
		}
	}
	projects, err := discovery.DeepProjectDiscoverer{}.Discover(root)
	if err != nil {
		return err
	}

	eg, nctx := errgroup.WithContext(ctx)
	eg.SetLimit(10)

	for _, p := range projects {
		eg.Go(initialization.NewProjectIniter(nctx, p).
			WithAssetDiscoverer(discovery.FlatAssetDiscoverer{}).
			PersistOnFinish().
			GetRunner(nil))
	}
	eg.Wait()

	return nil
}
