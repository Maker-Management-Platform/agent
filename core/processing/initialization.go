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

func EnqueueInitJob(asset *ProcessableAsset) {
	queue.Enqueue(asset)
}

func (pa *ProcessableAsset) JobAction() {
	log.Println("Initializing asset", pa.Name)
	var err error
	if _, err = database.GetAssetByProjectAndName(pa.Project.UUID, pa.Name); err != gorm.ErrRecordNotFound {
		log.Println("Asset already exists")
		return
	}
	pa.Asset, err = entities.NewProjectAsset2(pa.Name, pa.Label, pa.Project, pa.Origin)
	if err != nil {
		log.Println(err)
		return
	}
	err = processType(pa)
	if err != nil {
		log.Println(err)
		return
	}
	if pa.Asset.AssetType == "image" {
		if pa.Project.DefaultImageID == "" {
			pa.Project.DefaultImageID = pa.Asset.ID
			err = database.SetProjectDefaultImage(pa.Project.UUID, pa.Asset.ID)
			if err != nil {
				log.Println(err)
			}
		}
		if pa.Parent != nil {
			err = database.UpdateAssetImage(pa.Parent.ID, pa.Asset.ID)
			if err != nil {
				log.Println(err)
			}
		}
	}
	err = database.InsertAsset(pa.Asset)
	if err != nil {
		log.Println(err)
		return
	}
}

func (pa *ProcessableAsset) JobName() string {
	return fmt.Sprintf("Initialize %s", pa.Name)
}

func processType(pa *ProcessableAsset) error {
	var err error

	if t, ok := state.ExtensionProjectType[pa.Asset.Extension]; ok {
		pa.Asset.AssetType = t.Name
	}
	QueueEnrichmentJob(pa)

	return err
}

func (pa *ProcessableAsset) OnEnrichmentComplete(err error) {
	if err != nil {
		log.Println(err)
		return
	}

	if err = database.UpdateAssetProperties(pa.Asset.ID, pa.Asset.Properties); err != nil {
		log.Println(err)
		return
	}
}
