package processing

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/eduardooliveira/stLib/core/entities"
	"github.com/eduardooliveira/stLib/core/processing/discovery"
	"github.com/eduardooliveira/stLib/core/processing/initialization"
	"github.com/eduardooliveira/stLib/core/processing/types"
	"github.com/eduardooliveira/stLib/core/runtime"
	"github.com/eduardooliveira/stLib/core/utils"
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
	outs := make([]chan *types.ProcessableProject, 0)
	for _, p := range projects {
		out, runner := utils.Jobber(initialization.NewProjectIniter(p).
			WithContext(nctx).
			WithAssetDiscoverer(discovery.FlatAssetDiscoverer{}).
			PersistOnFinish().
			Init)
		eg.Go(runner)
		outs = append(outs, out)
	}
	eg.Wait()

	for _, out := range outs {
		for p := range out {
			log.Println(p)
		}
	}

	return nil
}
