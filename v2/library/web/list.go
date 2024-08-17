package web

import (
	"errors"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/a-h/templ"
	"github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/library/web/comp"
	"github.com/eduardooliveira/stLib/v2/utils"
	"github.com/eduardooliveira/stLib/v2/web"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type listInput struct {
	c     echo.Context
	r     *http.Request
	Asset *entities.Asset
}

func (h webHandler) list(in *listInput) (rtn templ.Component, pgModel *comp.PaginationModel, err error) {

	filter := entities.Asset{}
	/*err = (&echo.DefaultBinder{}).BindQueryParams(in.c, &filter)
	if err != nil {
		return nil, err
	}*/
	filter.ParentID = &in.Asset.ID
	if utils.VoZ(filter.Kind) == "all" {
		filter.Kind = nil
	}

	var page int

	if in.r.URL.Query().Get("page") != "" {
		page, _ = strconv.Atoi(in.r.URL.Query().Get("page"))
	}

	if page < 1 {
		page = 1
	}
	var pages int
	if in.Asset.ID != "" {
		pages, err = h.r.GetPagedNested(in.Asset, &filter, page-1, 20)
		if err != nil {
			return nil, nil, err
		}
	}

	return comp.List(comp.ListModel{
			Asset: in.Asset,
			Pagination: comp.PaginationModel{
				TotalPages:  pages,
				CurrentPage: page,
			},
		}), &comp.PaginationModel{
			TotalPages:  pages,
			CurrentPage: page,
		}, nil
}

func (h webHandler) listHandler(r *http.Request) web.ResponseModel {
	var err error
	var asset entities.Asset
	if r.URL.Query().Get("assetID") == "" {
		return web.ResponseModel{
			S:     http.StatusBadRequest,
			Error: errors.New("Asset ID is required"),
		}
	}
	asset, err = h.r.GetAsset(r.URL.Query().Get("assetID"), true)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return web.ResponseModel{
				S:     http.StatusNotFound,
				Error: err,
			}
		}
		return web.ResponseModel{
			S:     http.StatusInternalServerError,
			Error: err,
		}
	}

	err = h.r.LoadParents(&asset, 5, "ID", "Label")
	if err != nil {
		return web.ResponseModel{
			S:     http.StatusInternalServerError,
			Error: err,
		}
	}

	listComp, pgModel, err := h.list(&listInput{
		Asset: &asset,
		r:     r,
	})

	if err != nil {
		return web.ResponseModel{
			S:     http.StatusInternalServerError,
			Error: err,
		}
	}

	u, err := url.Parse(r.URL.String())
	if err != nil {
		return web.ResponseModel{
			S:     http.StatusInternalServerError,
			Error: err,
		}
	}

	u.Path = path.Join("/", "lib", asset.ID)
	q := u.Query()
	q.Del("assetID")
	u.RawQuery = q.Encode()

	pgModel.OOB = true
	return web.ResponseModel{
		S:          http.StatusOK,
		Component:  listComp,
		IsFragment: true,
		PushState:  u.String(),
		OOB: []templ.Component{
			comp.Pagination(*pgModel),
		},
	}
}
