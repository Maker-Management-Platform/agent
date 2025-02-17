package api

import (
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/library/libfs"
	"github.com/eduardooliveira/stLib/v2/utils"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func (h APIHandler) getFileHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("assetID")
	if id == "" {
		http.Error(w, "Asset ID is required", http.StatusBadRequest)
		return
	}
	asset, err := h.repo.GetAsset(id, false)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if chi.URLParam(r, "download") != "" {
		w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(*asset.Path))
	} else {
		w.Header().Set("Content-Disposition", "inline; filename="+filepath.Base(*asset.Path))
	}

	if shouldCacheFile(asset) {
		if err := h.repo.LoadTree(&asset, func(a *entities.Asset) bool { return false }); err != nil {
			h.log.Error("get file", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cacheFS, err := libfs.GetLibFS("cache")
		if err != nil {
			h.log.Error("get file", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cachePath := filepath.Join(utils.VoZ(asset.ParentID), filepath.Base(*asset.Path))
		_, err = fs.Stat(cacheFS, cachePath)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				h.log.Error("get file", "error", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				h.log.Info("create folder", "folder", filepath.Join(utils.VoZ(asset.ParentID)))
				target, err := cacheFS.Create(cachePath)
				if err != nil {
					h.log.Error("create file", "error", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return

				}
				defer target.Close()

				fs, err := libfs.GetAssetFS(r.Context(), asset)
				if err != nil {
					h.log.Error("get fs", "error", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				file, err := fs.Open(utils.VoZ(asset.Path))
				if err != nil {
					h.log.Error("open file", "error", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer file.Close()

				_, err = io.Copy(target, file)
				if err != nil {
					h.log.Error("copy file", "error", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

			}
		}

		http.ServeFileFS(w, r, cacheFS, cachePath)
		return
	}

	fs, err := libfs.GetAssetFS(r.Context(), asset)
	if err != nil {
		h.log.Error("get fs", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.ServeFileFS(w, r, fs, *asset.Path)
}

func shouldCacheFile(asset entities.Asset) bool {
	if asset.FSKind != "cache" && asset.FSKind != "generated" && asset.FSKind != "local" {
		return true
	}
	return false
}
