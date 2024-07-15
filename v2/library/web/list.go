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
	"github.com/eduardooliveira/stLib/v2/web"
	corecomp "github.com/eduardooliveira/stLib/v2/web/comp"
	"github.com/eduardooliveira/stLib/v2/web/mw"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type listInput struct {
	c     echo.Context
	Asset *entities.Asset
}

func (h webHandler) list(in *listInput) (rtn templ.Component, err error) {

	var page int

	if in.c.QueryParam("page") != "" {
		page, _ = strconv.Atoi(in.c.QueryParam("page"))
	}

	if page < 1 {
		page = 1
	}

	pages, err := h.r.GetPagedNested(in.Asset, page-1, 20)
	if err != nil {
		return nil, err
	}
	uh := mw.URLHelper(in.c).Clone()
	uh.SetTotalPages(pages)
	uh.SetPage(page)
	uh.WithPath(path.Join("/", "lib", in.Asset.ID, "list"))

	return comp.List(comp.ListModel{
		Asset: in.Asset,
		Pagination: corecomp.PaginationModel{
			TotalPages:  pages,
			CurrentPage: page,
			UH:          uh,
			Target:      ".asset-view .main",
		},
	}), nil
}

func (h webHandler) listHandler(c echo.Context) error {
	var err error
	var asset entities.Asset
	if c.Param("assetID") == "" {
		return web.Error(c, http.StatusBadRequest, "Asset ID is required")
	}
	asset, err = h.r.GetAsset(c.Param("assetID"), true)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return web.Error(c, http.StatusNotFound, "Asset not found")
		}
		return web.Error(c, http.StatusInternalServerError, err.Error())
	}

	err = h.r.LoadParents(&asset, 5, "ID", "Label")
	if err != nil {
		return web.Error(c, http.StatusInternalServerError, err.Error())
	}

	listComp, err := h.list(&listInput{
		c:     c,
		Asset: &asset,
	})

	if err != nil {
		return web.Error(c, http.StatusInternalServerError, err.Error())
	}

	u, err := url.Parse(c.Request().URL.String())
	if err != nil {
		return web.Error(c, http.StatusInternalServerError, err.Error())
	}

	return web.Render(web.ResponseModel{
		Ctx: c,
		S:   http.StatusOK,
		WrapperModel: corecomp.WrapperModel{
			Main: listComp,
		},
		IsFragment: true,
		PushState:  u.String(),
	})
}
