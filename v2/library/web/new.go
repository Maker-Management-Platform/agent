package web

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/eduardooliveira/stLib/v2/config"
	"github.com/eduardooliveira/stLib/v2/library/downloader"
	"github.com/eduardooliveira/stLib/v2/library/downloader/tools"
	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/library/web/comp"
	"github.com/eduardooliveira/stLib/v2/utils"
	"github.com/eduardooliveira/stLib/v2/web"
	"github.com/go-chi/chi/v5"
	"github.com/labstack/echo/v4"
)

type newAssetRequest struct {
	Mode     string `form:"mode"`
	ParentID string `form:"parentID"`
	Urls     string `form:"urls"`
	TempFile string `form:"tempFile"`
	Folder   string `form:"folder"`
}

func (h webHandler) newAsset(r *http.Request) web.ResponseModel {
	events := []string{}
	model := &comp.NewModel{
		TempFiles: []string{},
	}
	if r.Method == http.MethodPost {
		/*var req newAssetRequest
		err := c.Bind(&req)
		if err != nil {
			return web.Error(c, http.StatusBadRequest, err.Error())
		}

		if req.ParentID == "" {
			return web.Error(c, http.StatusBadRequest, "missing parentID")
		}

		parent, err := h.r.GetAsset(req.ParentID, false)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return web.Error(c, http.StatusNotFound, err.Error())
			}
			return web.Error(c, http.StatusInternalServerError, err.Error())
		}
		model.ParentID = req.ParentID

		switch req.Mode {
		case "import":
			if err := h.handleDownload(c, parent, req); err != nil {
				return err
			}
		case "upload":
			if err := h.handleUpload(c, parent, req); err != nil {
				return err
			}
		case "tempFiles":
			if err := h.handleTempFiles(c, parent, req); err != nil {
				return err
			}
		}
		events = append(events, "nested-assets-update")
		*/
	} else if r.Method == http.MethodGet {
		model.ParentID = chi.URLParam(r, "assetID")
	}
	entries, err := os.ReadDir(filepath.Join(config.Cfg.Core.DataFolder, "temp"))
	if err != nil {
		return web.ResponseModel{
			Error: err,
			S:     http.StatusInternalServerError,
		}
	}

	for _, e := range entries {
		model.TempFiles = append(model.TempFiles, e.Name())
	}
	return web.ResponseModel{
		S:          http.StatusOK,
		Component:  comp.New(model),
		Events:     events,
		IsFragment: true,
	}
}

func (h webHandler) handleDownload(c echo.Context, parent entities.Asset, req newAssetRequest) error {
	if req.Urls == "" {
		return web.Error(c, http.StatusBadRequest, "missing urls")
	}
	h.l.Info("new", "urls", req.Urls)

	err := downloader.Download(downloader.DownloadInput{
		Parent: parent,
		URL:    req.Urls,
		R:      h.r,
		P:      h.p,
	})
	if err != nil {
		return web.Error(c, http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (h webHandler) handleUpload(c echo.Context, parent entities.Asset, req newAssetRequest) error {
	form, err := c.MultipartForm()
	if err != nil {
		h.l.Error("new", "err", err)
		return err
	}
	files := form.File["files"]
	if len(files) == 0 && req.Folder == "" {
		return web.Error(c, http.StatusBadRequest, "no files or folder")
	}

	if req.Folder != "" {
		err = utils.CreateFolder(filepath.Join(*parent.Root, *parent.Path, req.Folder))
		if err != nil {
			return fmt.Errorf("creating folder: %v", err)
		}
		f := entities.NewAssetFromRootPath(*parent.Root, filepath.Join(*parent.Path, req.Folder), true, &parent)

		if err := h.r.SaveAsset(*f); err != nil {
			h.l.Error("new", "err", err)
			return err
		}
		parent = *f
	}

	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			h.l.Error("new", "err", err)
			return err
		}
		defer src.Close()
		err = tools.SaveFile(filepath.Join(*parent.Root, *parent.Path, file.Filename), src)
		if err != nil {
			h.l.Error("new", "err", err)
			return err
		}

		a := entities.NewAssetFromRootPath(*parent.Root, filepath.Join(*parent.Path, file.Filename), false, &parent)
		if err := h.r.SaveAsset(*a); err != nil {
			h.l.Error("new", "err", err)
			return err
		}
		if err := h.p.Process(a).Wait(); err != nil {
			h.l.Error("new", "err", err)
			return err
		}

	}

	return nil
}

func (h webHandler) handleTempFiles(c echo.Context, parent entities.Asset, req newAssetRequest) error {
	if req.TempFile == "" {
		return web.Error(c, http.StatusBadRequest, "missing Temp File")
	}
	err := utils.Move(filepath.Join(config.Cfg.Core.DataFolder, "temp", req.TempFile),
		filepath.Join(*parent.Root, *parent.Path, req.TempFile))
	if err != nil {
		return web.Error(c, http.StatusInternalServerError, err.Error())
	}
	a := entities.NewAssetFromRootPath(*parent.Root, filepath.Join(*parent.Path, req.TempFile), false, &parent)
	if err := h.r.SaveAsset(*a); err != nil {
		return web.Error(c, http.StatusInternalServerError, err.Error())
	}
	if err := h.p.Process(a).Wait(); err != nil {
		return web.Error(c, http.StatusInternalServerError, err.Error())
	}
	return nil
}
