package processing

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/eduardooliveira/stLib/core/data/database"
	"github.com/eduardooliveira/stLib/core/entities"
	"github.com/eduardooliveira/stLib/core/queue"
	"github.com/eduardooliveira/stLib/core/render"
)

type DiscoverableAsset struct {
	name    string
	path    string
	project *entities.Project
	parent  *entities.ProjectAsset
}

type initialize struct {
	da *DiscoverableAsset
}

func (i *initialize) Run() {
	if i.da.parent != nil {
		log.Println(i.da.parent.Name)
	}
	asset, err := entities.NewProjectAsset2(i.da.name, i.da.project)
	if err != nil {
		log.Println(err)
		return
	}
	err = processType(asset, i.da.project)
	if err != nil {
		log.Println(err)
		return
	}
	if asset.AssetType == "image" {
		if i.da.project.DefaultImageID == "" {
			i.da.project.DefaultImageID = asset.ID
			err = database.SetProjectDefaultImage(i.da.project.UUID, asset.ID)
			if err != nil {
				log.Println(err)
			}
		}
		if i.da.parent != nil {
			err = database.UpdateAssetImage(i.da.parent.ID, asset.ID)
			if err != nil {
				log.Println(err)
			}
		}
	}

	err = database.InsertAsset(asset)
	if err != nil {
		log.Println(err)
		return
	}
}

func (i *initialize) Name() string {
	return fmt.Sprintf("Initialize %s", i.da.name)
}

func Enqueue(asset *DiscoverableAsset) {
	queue.Enqueue(&initialize{da: asset})
}

func processType(asset *entities.ProjectAsset, project *entities.Project) error {
	var err error
	if slices.Contains(entities.ModelExtensions, strings.ToLower(asset.Extension)) {
		asset.AssetType = entities.ProjectModelType
		asset.Model, err = entities.NewProjectModel2(asset, project)
		if err == nil {
			render.QueueJob(&renderableAsset{asset: asset, project: project})
		}
	} else if slices.Contains(entities.ImageExtensions, strings.ToLower(asset.Extension)) {
		asset.ProjectImage, err = entities.NewProjectImage2(asset, project)
	} else if slices.Contains(entities.SliceExtensions, strings.ToLower(asset.Extension)) {
		asset.Slice, err = entities.NewProjectSlice2(asset, project)
	} else {
		asset.AssetType = entities.ProjectFileType
		asset.ProjectFile, err = entities.NewProjectFile2(asset, project)
	}
	for _, ext := range entities.GeneratedExtensions {
		if strings.HasSuffix(asset.Name, ext) {
			asset.Generated = true
		}
	}
	return err
}

type renderableAsset struct {
	asset   *entities.ProjectAsset
	project *entities.Project
}

func (r *renderableAsset) Name() string {
	return r.asset.Name
}

func (r *renderableAsset) Project() *entities.Project {
	return r.project
}
func (r *renderableAsset) Asset() *entities.ProjectAsset {
	return r.asset
}

func (r *renderableAsset) OnComplete(path string, err error) {
	Enqueue(&DiscoverableAsset{
		name:    path,
		path:    path,
		project: r.project,
		parent:  r.asset,
	})
}
