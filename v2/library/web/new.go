package web

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/eduardooliveira/stLib/v2/library/downloader"
	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/library/libfs"
	"github.com/eduardooliveira/stLib/v2/library/web/comp"
	"github.com/eduardooliveira/stLib/v2/web"
	"github.com/ggicci/httpin"
	"gorm.io/gorm"
)

type newAssetRequest struct {
	Mode     string         `in:"form=mode"`
	ParentID string         `in:"form=parentID"`
	Urls     string         `in:"form=urls"`
	TempFile string         `in:"form=tempFile"`
	Folder   string         `in:"form=folder"`
	Files    []*httpin.File `in:"form=files"`
}

func (h webHandler) newAssetHandler(r *http.Request) web.ResponseModel {
	events := []string{}
	model := &comp.NewModel{
		TempFiles: []string{},
	}
	if r.Method == http.MethodPost {
		req := r.Context().Value(httpin.Input).(*newAssetRequest)

		if req.ParentID == "" {
			return web.ResponseModel{
				Error:  errors.New("missing parentID"),
				Status: http.StatusBadRequest,
			}
		}

		parent, err := h.r.GetAsset(req.ParentID, false)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return web.ResponseModel{
					Error:  err,
					Status: http.StatusNotFound,
				}
			}
			return web.ResponseModel{
				Error:  err,
				Status: http.StatusInternalServerError,
			}
		}
		model.ParentID = req.ParentID

		reqHandlers := map[string]func(r *http.Request, parent entities.Asset, req newAssetRequest) (error, int){
			"import":    h.handleDownload,
			"upload":    h.handleUpload,
			"tempFiles": h.handleTempFiles,
		}
		if req.Mode == "" {
			return web.ResponseModel{
				Error:  errors.New("missing mode"),
				Status: http.StatusBadRequest,
			}
		}
		if _, ok := reqHandlers[req.Mode]; !ok {
			return web.ResponseModel{
				Error:  errors.New("invalid mode"),
				Status: http.StatusBadRequest,
			}
		}
		if err, s := reqHandlers[req.Mode](r, parent, *req); err != nil {
			return web.ResponseModel{
				Error:  err,
				Status: s,
			}
		}

		events = append(events, "nested-assets-update")

	} else if r.Method == http.MethodGet {
		h.l.Info("new", "GET", r.URL.Query().Get("assetID"))
		model.ParentID = r.URL.Query().Get("assetID")
	}
	//TODO: implement uploadFS
	/*entries, err := os.ReadDir(filepath.Join(config.Cfg.Core.DataFolder, "temp"))
	if err != nil {
		return web.ResponseModel{
			Error:  err,
			Status: http.StatusInternalServerError,
		}
	}

	for _, e := range entries {
		model.TempFiles = append(model.TempFiles, e.Name())
	}*/
	return web.ResponseModel{
		Status:     http.StatusOK,
		Component:  comp.New(model),
		Events:     events,
		IsFragment: true,
	}
}

func (h webHandler) handleDownload(r *http.Request, parent entities.Asset, req newAssetRequest) (error, int) {
	if req.Urls == "" {
		return errors.New("missing urls"), http.StatusBadRequest
	}
	h.l.Info("new", "urls", req.Urls)

	err := downloader.Download(downloader.DownloadInput{
		Parent:    parent,
		URL:       req.Urls,
		Repo:      h.r,
		Processor: h.p,
	})
	if err != nil {
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}

func (h webHandler) handleUpload(r *http.Request, parent entities.Asset, req newAssetRequest) (error, int) {
	var err error

	files := req.Files
	if len(files) == 0 && req.Folder == "" {
		return errors.New("no files or folder"), http.StatusBadRequest
	}

	fSys, err := libfs.GetFS(parent.FSKind, parent.FSName, *parent.Root)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	if req.Folder != "" {
		err = fSys.Mkdir(filepath.Join(*parent.Path, req.Folder))
		if err != nil {
			return fmt.Errorf("creating folder: %v", err), http.StatusInternalServerError
		}
		f := entities.NewAsset(fSys, filepath.Join(*parent.Path, req.Folder), true, &parent)

		if err := h.r.SaveAsset(*f); err != nil {
			h.l.Error("new", "err", err)
			return err, http.StatusInternalServerError
		}
		parent = *f
	}

	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			h.l.Error("new", "err", err)
			return err, http.StatusInternalServerError
		}
		defer src.Close()
		out, err := fSys.Create(filepath.Join(*parent.Path, file.Filename()))
		defer out.Close()
		if err != nil {
			h.l.Error("new", "err", err)
			return err, http.StatusInternalServerError
		}

		_, err = io.Copy(out, src)
		if err != nil {
			h.l.Error("new", "err", err)
			return err, http.StatusInternalServerError
		}

		a := entities.NewAsset(fSys, filepath.Join(*parent.Path, file.Filename()), false, &parent)
		if err := h.r.SaveAsset(*a); err != nil {
			h.l.Error("new", "err", err)
			return err, http.StatusInternalServerError
		}
		if err := h.p.Process(a).Wait(); err != nil {
			h.l.Error("new", "err", err)
			return err, http.StatusInternalServerError
		}

	}

	return nil, http.StatusOK
}

func (h webHandler) handleTempFiles(r *http.Request, parent entities.Asset, req newAssetRequest) (error, int) {
	/*if req.TempFile == "" {
		return errors.New("missing Temp File"), http.StatusBadRequest
	}
	err := utils.Move(filepath.Join(config.Cfg.Core.DataFolder, "temp", req.TempFile),
		filepath.Join(*parent.Root, *parent.Path, req.TempFile))
	if err != nil {
		return err, http.StatusInternalServerError
	}
	a := entities.NewAssetFromRootPath(*parent.Root, filepath.Join(*parent.Path, req.TempFile), false, &parent)
	if err := h.r.SaveAsset(*a); err != nil {
		return err, http.StatusInternalServerError
	}
	if err := h.p.Process(a).Wait(); err != nil {
		return err, http.StatusInternalServerError
	}*/
	return nil, http.StatusOK
}
