package initialization

import (
	"context"
	"log"

	"github.com/eduardooliveira/stLib/core/data/database"
	"github.com/eduardooliveira/stLib/core/entities"
	"github.com/eduardooliveira/stLib/core/processing/enrichment"
	"github.com/eduardooliveira/stLib/core/processing/types"
	"github.com/eduardooliveira/stLib/core/state"
)

type AssetIniter struct {
	ctx context.Context
	pa  *types.ProcessableAsset
	out chan<- *entities.ProjectAsset
}

func NewAssetIniter(ctx context.Context, pa *types.ProcessableAsset) *AssetIniter {
	return &AssetIniter{
		ctx: ctx,
		pa:  pa,
	}
}

func (ai *AssetIniter) GetRunner(out chan *entities.ProjectAsset) func() error {
	ai.out = out
	return ai.run
}

func (ai *AssetIniter) run() error {
	defer ai.close()

	if a, err := database.GetAssetByProjectAndName(ai.pa.Project.UUID, ai.pa.Name); err == nil && a.ID != "" {
		ai.pa.Asset = a
	} else {
		ai.pa.Asset, err = entities.NewProjectAsset2(ai.pa.Name, ai.pa.Label, ai.pa.Project, ai.pa.Origin)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	if err := ai.processType(); err != nil {
		log.Println(err)
		return err
	}

	nestedAssets, err := enrichment.EnrichAsset(ai.ctx, ai.pa)
	if err != nil {
		log.Println(err)
	}

	for _, nestedAsset := range nestedAssets {

		if err := NewAssetIniter(ai.ctx, nestedAsset).GetRunner(nil)(); err != nil {
			log.Println(err)
		}
		if nestedAsset.Asset.AssetType == "image" {
			ai.pa.Asset.ImageID = nestedAsset.Asset.ID
		}
		ai.out <- nestedAsset.Asset
	}

	if err := database.SaveAsset(ai.pa.Asset); err != nil {
		log.Println(err)
		return err
	}
	if ai.out != nil {
		ai.out <- ai.pa.Asset
	}
	return nil
}

func (ai *AssetIniter) close() {
	if ai.out != nil {
		close(ai.out)
	}
}

func (ai *AssetIniter) processType() error {
	if t, ok := state.ExtensionProjectType[ai.pa.Asset.Extension]; ok {
		ai.pa.Asset.AssetType = t.Name
	} else {
		ai.pa.Asset.AssetType = "other"
	}

	if ai.pa.Asset.AssetType == "image" {
		ai.pa.Asset.ImageID = ai.pa.Asset.ID
	}
	return nil
}
