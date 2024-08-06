package web

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/eduardooliveira/stLib/v2/config"
	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/library/sys"
	"github.com/eduardooliveira/stLib/v2/utils"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func (h webHandler) getFileHandlerChi(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "assetID")
	if id == "" {
		http.Error(w, "Asset ID is required", http.StatusBadRequest)
		return
	}
	asset, err := h.r.GetAsset(id, false)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if utils.VoZ(asset.NodeKind) == entities.NodeKindBundled {
		if err := h.r.LoadParents(&asset, 1); err != nil {
			h.l.Error("get file", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err := os.Stat(filepath.Join(config.Cfg.Core.DataFolder, "temp", utils.VoZ(asset.ParentID), filepath.Base(*asset.Path)))
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				h.l.Error("get file", "error", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				h.l.Info("create folder", "folder", filepath.Join(config.Cfg.Core.DataFolder, "temp", utils.VoZ(asset.ParentID)))
				err = utils.CreateFolder(filepath.Join(config.Cfg.Core.DataFolder, "temp", utils.VoZ(asset.ParentID)))
				if err != nil {
					h.l.Error("create folder", "error", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				fs, err := sys.GetFS(utils.VoZ(asset.FSKind), utils.VoZ(asset.FSName), *asset.Root)
				if err != nil {
					h.l.Error("get fs", "error", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				file, err := fs.Open(utils.VoZ(asset.Path))
				if err != nil {
					h.l.Error("open file", "error", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer file.Close()

				data, err := io.ReadAll(file)
				if err != nil {
					h.l.Error("read file", "error", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				err = os.WriteFile(filepath.Join(config.Cfg.Core.DataFolder, "temp", utils.VoZ(asset.ParentID), filepath.Base(*asset.Path)), data, 0644)
				if err != nil {
					h.l.Error("write file", "error", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

			}
		}

	}
	if chi.URLParam(r, "download") != "" {
		w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(*asset.Path))
	} else {
		w.Header().Set("Content-Disposition", "inline; filename="+filepath.Base(*asset.Path))
	}

	if utils.VoZ(asset.NodeKind) == entities.NodeKindBundled {
		http.ServeFile(w, r, filepath.Join(config.Cfg.Core.DataFolder, "temp", utils.VoZ(asset.ParentID), filepath.Base(*asset.Path)))
		return
	}
	http.ServeFile(w, r, filepath.Join(*asset.Root, *asset.Path))
}
