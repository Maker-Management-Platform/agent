package processing

import (
	"fmt"
	"log"

	"github.com/eduardooliveira/stLib/core/data/database"
	"github.com/eduardooliveira/stLib/core/entities"
	"github.com/eduardooliveira/stLib/core/queue"
	"github.com/eduardooliveira/stLib/core/state"
	"gorm.io/gorm"
)

type DiscoverableAsset struct {
	Name    string
	Path    string
	Project *entities.Project
	Parent  *entities.ProjectAsset
}

func EnqueueInitJob(asset *processableAsset) {
	queue.Enqueue(asset)
}

func (pa *processableAsset) JobAction() {
	log.Println("Initializing asset", pa.name)
	var err error
	if _, err = database.GetAssetByProjectAndName(pa.project.UUID, pa.name); err != gorm.ErrRecordNotFound {
		log.Println("Asset already exists")
		return
	}
	pa.asset, err = entities.NewProjectAsset2(pa.name, pa.label, pa.project, pa.origin)
	if err != nil {
		log.Println(err)
		return
	}
	err = processType(pa)
	if err != nil {
		log.Println(err)
		return
	}
	if pa.asset.AssetType == "image" {
		if pa.project.DefaultImageID == "" {
			pa.project.DefaultImageID = pa.asset.ID
			err = database.SetProjectDefaultImage(pa.project.UUID, pa.asset.ID)
			if err != nil {
				log.Println(err)
			}
		}
		if pa.parent != nil {
			err = database.UpdateAssetImage(pa.parent.ID, pa.asset.ID)
			if err != nil {
				log.Println(err)
			}
		}
	}
	err = database.InsertAsset(pa.asset)
	if err != nil {
		log.Println(err)
		return
	}
}

func (pa *processableAsset) JobName() string {
	return fmt.Sprintf("Initialize %s", pa.name)
}

func processType(pa *processableAsset) error {
	var err error

	if t, ok := state.ExtensionProjectType[pa.asset.Extension]; ok {
		pa.asset.AssetType = t.Name
	}
	QueueEnrichmentJob(pa)

	return err
}

func (pa *processableAsset) OnEnrichmentComplete(err error) {
	if err != nil {
		log.Println(err)
		return
	}

	if err = database.UpdateAssetProperties(pa.asset.ID, pa.asset.Properties); err != nil {
		log.Println(err)
		return
	}
}
