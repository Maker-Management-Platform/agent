package helpers

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
)

type URLHelper struct {
	URL *url.URL
}

func NewURLHelper(req *http.Request) *URLHelper {
	return &URLHelper{
		URL: req.URL,
	}
}

func (u *URLHelper) SetParam(key, value string) {
	q := u.URL.Query()
	q.Set(key, value)
	u.URL.RawQuery = q.Encode()
}

func (u *URLHelper) SetPage(page int) {
	u.SetParam("page", fmt.Sprintf("%d", page))
}

func (u *URLHelper) SetTotalPages(perPage int) {
	u.SetParam("totalPages", fmt.Sprintf("%d", perPage))
}

func (u *URLHelper) WithPage(page int) *URLHelper {
	u.SetPage(page)
	return u
}

func (u *URLHelper) WithPagePlus(plus int) *URLHelper {
	p := u.URL.Query().Get("page")
	if p == "" {
		p = "0"
	}
	pInt, err := strconv.Atoi(p)
	if err != nil {
		slog.Error("Error converting page to int")
	}

	u.SetParam("page", fmt.Sprintf("%d", plus+pInt))
	return u
}

func (u *URLHelper) WithPath(path string) *URLHelper {
	u.URL.Path = path
	return u
}

func (u URLHelper) WithQuery(key, value string) URLHelper {
	u.SetParam(key, value)
	return u
}

func (u URLHelper) GetQuery() string {
	return u.URL.RawQuery
}
func (u URLHelper) GetURL() string {
	return u.URL.String()
}

func (u *URLHelper) Clone() *URLHelper {
	nu := *u.URL
	return &URLHelper{
		URL: &nu,
	}
}
