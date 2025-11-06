package api

import (
	"encoding/json"
	"net/http"

	"github.com/eduardooliveira/stLib/v2/api"
	"github.com/eduardooliveira/stLib/v2/utils"
)

type AssetPatch struct {
	Label       utils.Option[string]   `json:"label"`
	Description utils.Option[string]   `json:"description"`
	Tags        utils.Option[[]string] `json:"tags"`
}

func (a *APIHandler) patchHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("assetID")
	patch := AssetPatch{}

	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		api.Error(w, r, http.StatusBadRequest, err)
		return
	}

	asset, err := a.repo.GetAsset(id, false)
	if err != nil {
		api.Error(w, r, http.StatusNotFound, err)
		return
	}

	if patch.Label.Valid {
		asset.Label = &patch.Label.Value
	}
	if patch.Description.Valid {
		asset.Description = &patch.Description.Value
	}
	if patch.Tags.Valid {
		a.log.Warn("Tags are not implemented yet")
	}

	if err := a.repo.UpdateAsset(&asset); err != nil {
		api.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	api.JSON(w, r, http.StatusOK, asset)
}
