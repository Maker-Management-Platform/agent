package klipper

import (
	"encoding/json"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/duke-git/lancet/v2/maputil"
	"github.com/eduardooliveira/stLib/core/events"
	"github.com/gorilla/websocket"
)

func (kp KlipperPrinter) listen() {
	for {
		if err := kp.connect(); err != nil {
			log.Println(err)
			time.Sleep(30 * time.Second)
			continue
		}
		for {
			_, message, err := kp.ws.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			kpStatusString := string(message)

			statusDelta := make([]*events.Message, 0)
			if strings.Contains(kpStatusString, "notify_proc_stat_update") {
				continue
			} else if strings.Contains(kpStatusString, "notify_status_update") {
				statusDelta = kp.handleNotifyStatusUpdate(message)
			} else if strings.Contains(kpStatusString, "result") {
				statusDelta = kp.handleQuery(message)
			}
			broadcast[[]*events.Message](kp.deltaListeners, statusDelta)
		}
	}
}

func (kp KlipperPrinter) handleNotifyStatusUpdate(message []byte) []*events.Message {
	var kpStatusUpdate *statusUpdate
	err := json.Unmarshal(message, &kpStatusUpdate)
	if err != nil {
		log.Println(err)
		return nil
	}
	status := make(map[string]*events.Message, 0)
	for _, p := range kpStatusUpdate.Params {
		if param, ok := p.(map[string]any); ok {
			for k, v := range param {
				addToStatus(k, v.(map[string]any), status)
			}
		}
	}
	return maputil.Values(status)
}

func (kp KlipperPrinter) handleQuery(message []byte) []*events.Message {
	var pkResult *result
	err := json.Unmarshal(message, &pkResult)
	if err != nil {
		log.Println(err)
		return nil
	}

	status := make(map[string]*events.Message, 0)
	for k, v := range pkResult.Result.Status {
		addToStatus(k, v, status)
	}
	return maputil.Values(status)
}

func (kp KlipperPrinter) connect() error {
	log.Println("klipper connecting to ", kp.config.Address)
	u, err := url.Parse(kp.config.Address)
	if err != nil {
		return err
	}

	u.Scheme = "ws"
	u.Path = "/websocket"

	kp.ws, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	return kp.ws.WriteMessage(websocket.TextMessage, []byte("{\"jsonrpc\":\"2.0\",\"method\":\"printer.objects.subscribe\",\"params\":{\"objects\":{\"heaters\":null,\"idle_timeout\":null,\"print_stats\":null,\"display_status\":null,\"heater_bed\":null,\"fan\":null,\"heater_fan toolhead_cooling_fan\":null,\"extruder\":null}},\"id\":1}"))
}
