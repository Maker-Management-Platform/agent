package initialization

import (
	"context"
	"log"

	"github.com/eduardooliveira/stLib/core/data/database"
	"github.com/eduardooliveira/stLib/core/entities"
	"github.com/eduardooliveira/stLib/core/processing/discovery"
	"github.com/eduardooliveira/stLib/core/processing/types"
	"github.com/eduardooliveira/stLib/core/utils"
)

type ProjectIniter struct {
	ctx                context.Context
	processableProject types.ProcessableProject
	assetDiscoverer    discovery.AssetDiscoverer
	project            entities.Project
	persistOnFinish    bool
	out                chan<- *entities.Project
}

func NewProjectIniter(ctx context.Context, processableProject types.ProcessableProject) *ProjectIniter {
	return &ProjectIniter{
		ctx:                ctx,
		processableProject: processableProject,
	}
}

func (pd *ProjectIniter) WithProject(project entities.Project) *ProjectIniter {
	pd.project = project
	return pd
}

func (pd *ProjectIniter) WithAssetDiscoverer(ad discovery.AssetDiscoverer) *ProjectIniter {
	pd.assetDiscoverer = ad
	return pd
}

func (pd *ProjectIniter) PersistOnFinish() *ProjectIniter {
	pd.persistOnFinish = true
	return pd
}

func (pd *ProjectIniter) GetRunner(out chan<- *entities.Project) func() error {
	pd.out = out
	return pd.run
}

func (pd *ProjectIniter) run() error {
	defer pd.close()
	pd.LoadProject()

	if pd.assetDiscoverer != nil {
		assets, err := pd.assetDiscoverer.Discover(pd.processableProject.Path)
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println(len(assets))

		outChans := make([]<-chan *entities.ProjectAsset, 0)
		for _, a := range assets {
			a.Project = &pd.project //TODO: make asset self-contained
			ai := NewAssetIniter(pd.ctx, a)
			c := make(chan *entities.ProjectAsset)
			go ai.GetRunner(c)()
			outChans = append(outChans, c)
		}

		out := utils.MergeWait(outChans...)

		for a := range out {
			log.Println(a)
			if a.AssetType == "image" {
				pd.project.DefaultImageID = a.ID
			}
		}
	}

	if pd.persistOnFinish {
		if err := database.InsertProject(&pd.project); err != nil {
			log.Println(err)
			return err
		}
	}

	if pd.out != nil {
		pd.out <- &pd.project
	}

	return nil
}

func (pd *ProjectIniter) Project() *entities.Project {
	return &pd.project
}

func (pd *ProjectIniter) close() {
	if pd.out != nil {
		close(pd.out)
	}
}

func (pd *ProjectIniter) LoadProject() {
	if pd.project.UUID != "" {
		return
	}

	pd.project = *entities.NewProjectFromPath(pd.processableProject.Path)
	if p, err := database.GetProjectByPathAndName(pd.project.Path, pd.project.Name); err == nil {
		pd.project = *p
	} else {
		pd.project.Tags = append(pd.project.Tags, entities.StringsToTags(utils.PathToTags(pd.project.Path))...)
	}
}
