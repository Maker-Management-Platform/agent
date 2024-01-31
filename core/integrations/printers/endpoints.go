package printers

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/eduardooliveira/stLib/core/data/database"
	"github.com/eduardooliveira/stLib/core/integrations/klipper"
	"github.com/eduardooliveira/stLib/core/models"
	"github.com/eduardooliveira/stLib/core/state"
	"github.com/labstack/echo/v4"
)

func index(c echo.Context) error {
	rtn := make([]*models.Printer, 0)
	for _, p := range state.Printers {
		rtn = append(rtn, p)
	}
	return c.JSON(http.StatusOK, rtn)
}

func show(c echo.Context) error {
	uuid := c.Param("uuid")
	printer, ok := state.Printers[uuid]

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, printer)
}

func deleteHandler(c echo.Context) error {
	uuid := c.Param("uuid")
	printer, ok := state.Printers[uuid]

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	delete(state.Printers, printer.UUID)

	err := state.PersistPrinters()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, printer)
}

func sendHandler(c echo.Context) error {
	uuid := c.Param("uuid")
	printer, ok := state.Printers[uuid]

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	id := c.Param("id")
	asset, err := database.GetAsset(id)

	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if printer.Type == "klipper" {
		err = klipper.UploadFile(printer, asset)
	}
	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func new(c echo.Context) error {

	pPrinter := &models.Printer{}
	err := c.Bind(pPrinter)

	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	printer := models.NewPrinter()
	printer.Name = pPrinter.Name
	printer.Address = pPrinter.Address
	printer.Type = pPrinter.Type

	state.Printers[printer.UUID] = printer
	state.PersistPrinters()

	return c.JSON(http.StatusCreated, state.Printers[printer.UUID])
}

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

func edit(c echo.Context) error {
	pPrinter := &models.Printer{}

	if err := c.Bind(pPrinter); err != nil {
		log.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}

	if pPrinter.UUID != c.Param("uuid") {
		return c.NoContent(http.StatusBadRequest)
	}

	printer, ok := state.Printers[pPrinter.UUID]

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	printer.Name = pPrinter.Name
	printer.Address = pPrinter.Address
	printer.Type = pPrinter.Type

	state.Printers[printer.UUID] = printer
	state.PersistPrinters()

	return c.JSON(http.StatusCreated, state.Printers[printer.UUID])
}

func testConnection(c echo.Context) error {

	pPrinter := &models.Printer{}
	err := c.Bind(pPrinter)

	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if pPrinter.Type == "klipper" {
		err = klipper.ConnectionStatus(pPrinter)
		log.Println(err)
		return c.JSON(http.StatusOK, pPrinter)
	}

	return c.NoContent(http.StatusInternalServerError)
}

func statusHandler(c echo.Context) error {

	uuid := c.Param("uuid")
	printer, ok := state.Printers[uuid]

	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, errors.New("printer not found").Error())
	}
	sm, err := GetStateManager(printer)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	stateChan := sm.Subscribe(c.Request().Context())

	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().WriteHeader(http.StatusOK)

	enc := json.NewEncoder(c.Response())
	for s := range stateChan {
		c.Response().Write([]byte("data: "))
		if err := enc.Encode(s); err != nil {
			return err
		}
		c.Response().Write([]byte("\n\n"))
		c.Response().Flush()
	}

	return c.JSON(http.StatusOK, printer)
}
