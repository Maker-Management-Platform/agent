package initialization

import (
	"fmt"
	"log"
	"regexp"
	"slices"
	"strings"

	"github.com/eduardooliveira/stLib/core/data/database"
	"github.com/eduardooliveira/stLib/core/entities"
	"github.com/eduardooliveira/stLib/core/processing/enrichment"
	"github.com/eduardooliveira/stLib/core/queue"
	"github.com/eduardooliveira/stLib/core/state"
)

type DiscoverableAsset struct {
	Name       string
	Path       string
	Project    *entities.Project
	Parent     *entities.ProjectAsset
	SkipInsert bool
}

type initialize struct {
	da *DiscoverableAsset
}

var parentRegex = regexp.MustCompile(`^(?P<parent>.*)\.(?:thumb|render)`)

func (i *initialize) Run() {
	if i.da.Parent != nil {
		log.Println(i.da.Parent.Name)
	}
	asset, err := entities.NewProjectAsset2(i.da.Name, i.da.Project)
	if err != nil {
		log.Println(err)
		return
	}
	err = processType(asset, i.da.Project)
	if err != nil {
		log.Println(err)
		return
	}
	if asset.AssetType == "image" {
		if i.da.Project.DefaultImageID == "" {
			i.da.Project.DefaultImageID = asset.ID
			err = database.SetProjectDefaultImage(i.da.Project.UUID, asset.ID)
			if err != nil {
				log.Println(err)
			}
		}
		if i.da.Parent == nil {
			match := parentRegex.FindStringSubmatch(i.da.Name)
			if len(match) == 2 {
				i.da.Parent, err = database.GetAssetByProjectAndName(i.da.Project.UUID, match[1])
				if err != nil {
					log.Println(err)
				}
			}
		}
		if i.da.Parent != nil {
			err = database.UpdateAssetImage(i.da.Parent.ID, asset.ID)
			if err != nil {
				log.Println(err)
			}
		}
	}
	if !i.da.SkipInsert {
		err = database.InsertAsset(asset)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (i *initialize) Name() string {
	return fmt.Sprintf("Initialize %s", i.da.Name)
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
			enrichment.QueueJob(&renderableAsset{asset: asset, project: project})
		}
	} else if slices.Contains(entities.ImageExtensions, strings.ToLower(asset.Extension)) {
		asset.ProjectImage, err = entities.NewProjectImage2(asset, project)
	} else if slices.Contains(entities.SliceExtensions, strings.ToLower(asset.Extension)) {
		asset.Slice, err = entities.NewProjectSlice2(asset, project)
		if err == nil {
			enrichment.QueueJob(&renderableAsset{asset: asset, project: project})
		}
	} else {
		asset.AssetType = entities.ProjectFileType
		asset.ProjectFile, err = entities.NewProjectFile2(asset, project)
	}
	for _, ext := range entities.GeneratedExtensions {
		if strings.HasSuffix(asset.Name, ext) {
			asset.Generated = true
		}
	}
	if t, ok := state.ExtensionProjectType[asset.Extension]; ok {
		asset.AssetType = t.Name
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
		Name:       path,
		Path:       path,
		Project:    r.project,
		Parent:     r.asset,
		SkipInsert: err != nil && err.Error() == "already exists",
	})
}
