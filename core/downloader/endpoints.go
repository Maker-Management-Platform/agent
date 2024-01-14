package downloader

import (
	"log"
	"net/http"
	"strings"

	"github.com/eduardooliveira/stLib/core/downloader/makerworld"
	"github.com/eduardooliveira/stLib/core/downloader/thingiverse"
	"github.com/labstack/echo/v4"
)

type urls struct {
	Url     string    `json:"url"`
	Urls    []string  `json:"urls"`
	Cookies []*cookie `json:"cookies"`
}
type cookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func fetch(c echo.Context) error {

	payload := &urls{
		Cookies: make([]*cookie, 0),
	}
	if err := c.Bind(payload); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if payload.Url != "" {
		payload.Urls = append(payload.Urls, strings.Split(payload.Url, ",")...)
	}

	for _, url := range payload.Urls {
		log.Println(url)
		if strings.Contains(url, "thingiverse.com") || strings.Contains(url, "thing:") {
			err := thingiverse.Fetch(url)
			if err != nil {
				log.Println(err)
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
		} else if strings.Contains(url, "makerworld.com") {
			cookies := make([]*http.Cookie, 0)

			for _, c := range payload.Cookies {
				cookies = append(cookies, &http.Cookie{
					Name:  c.Name,
					Value: c.Value,
					Path:  "/",
				})
			}
			agent := c.Request().UserAgent()
			err := makerworld.Fetch(url, cookies, agent)
			if err != nil {
				log.Println(err)
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
		}
	}

	return c.NoContent(http.StatusOK)
}
