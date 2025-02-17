package api

import (
	"net/http"

	"github.com/eduardooliveira/stLib/v2/api"
	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/utils"
)

func (a *APIHandler) indexHandler(w http.ResponseWriter, r *http.Request) {

	if r.PathValue("assetID") != "" {
		asset, err := a.repo.GetAsset(r.PathValue("assetID"), true)
		if err != nil {
			api.Error(w, r, http.StatusInternalServerError, err)
			return
		}
		err = a.repo.LoadParents(&asset, 5, "ID", "Label")
		if err != nil {
			api.Error(w, r, http.StatusInternalServerError, err)
			return
		}
		api.JSON(w, r, http.StatusOK, asset)
		return
	}

	roots, err := a.repo.GetAssetRoots(false)
	if err != nil {
		api.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	api.JSON(w, r, http.StatusOK, entities.Asset{
		Label:        utils.Ptr("Libraries"),
		NestedAssets: roots,
	})
}
