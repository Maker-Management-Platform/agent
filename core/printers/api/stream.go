package api

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/eduardooliveira/stLib/core/state"
	"github.com/labstack/echo/v4"
)

func stream(c echo.Context) error {
	uuid := c.Param("uuid")
	printer, ok := state.Printers[uuid]

	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, errors.New("printer not found").Error())
	}

	if printer.CameraUrl == "" {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("no configured camera").Error())
	}

	cameraUrl, err := url.Parse(printer.CameraUrl)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid camera url").Error())
	}

	req := &http.Request{
		Method: "GET",
		URL:    cameraUrl,
		/*Header: http.Header{
			"Authorization": []string{"Bearer " + runtime.Cfg.ThingiverseToken},
		},*/
	}
	httpClient := &http.Client{}

	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	for h, vs := range res.Header {
		for _, v := range vs {
			c.Response().Writer.Header().Add(h, v)
		}
	}

	r := bufio.NewReader(res.Body)
	buf := make([]byte, 0, 64*1024)
	for {
		select {
		case <-c.Request().Context().Done():
			return c.NoContent(http.StatusOK)
		default:
			n, err := io.ReadFull(r, buf[:cap(buf)])
			buf = buf[:n]
			if err != nil {
				if err == io.EOF {
					return echo.NewHTTPError(http.StatusBadRequest, err.Error())
				}
				if err != io.ErrUnexpectedEOF {
					return echo.NewHTTPError(http.StatusBadRequest, err.Error())
				}
			}
			c.Response().Write(buf)
			c.Response().Flush()
		}

	}
}
