package klipper

import (
	"encoding/json"
	"log"
	"net/url"
	"strings"

	"github.com/eduardooliveira/stLib/core/models"
	"github.com/gorilla/websocket"
)

type statePublisher struct {
	printer    *KlipperPrinter
	out        chan *models.PrinterStatus
	onNewSub   chan struct{}
	onShutdown chan struct{}
	conn       *websocket.Conn
}

func GetStatePublisher(printer *models.Printer) *statePublisher {
	kp := &KlipperPrinter{printer}
	return &statePublisher{
		printer:    kp,
		out:        make(chan *models.PrinterStatus),
		onNewSub:   make(chan struct{}),
		onShutdown: make(chan struct{}),
	}
}

func (p *statePublisher) Start() error {
	u, err := url.Parse(p.printer.Address)
	if err != nil {
		log.Println(err)
		return err
	}

	u.Scheme = "ws"
	u.Path = "/websocket"

	log.Printf("connecting to %s", u.String())
	p.conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println(err)
		return err
	}
	p.conn.WriteMessage(websocket.TextMessage, []byte("{\"jsonrpc\":\"2.0\",\"method\":\"printer.objects.subscribe\",\"params\":{\"objects\":{\"heaters\":null,\"idle_timeout\":null,\"print_stats\":null,\"display_status\":null,\"heater_bed\":null,\"fan\":null,\"heater_fan toolhead_cooling_fan\":null,\"extruder\":null}},\"id\":1}"))

	return nil
}

func (p *statePublisher) Produce() ([]*models.PrinterStatus, error) {
	select {
	case <-p.onNewSub:
		p.conn.WriteMessage(websocket.TextMessage, []byte("{\"jsonrpc\": \"2.0\",\"method\": \"printer.objects.query\",\"params\": {\"objects\": {\"extruder\": null,\"heater_bed\": null, \"display_status\": null}},\"id\": 2}"))
	default:
		_, message, err := p.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			//TODO: implement reconnect
			return nil, err
		}

		kpStatusString := string(message)

		if strings.Contains(kpStatusString, "notify_proc_stat_update") {
			return nil, nil
		}

		//log.Println(p.printer.Name, "status update:", kpStatusString)

		if strings.Contains(kpStatusString, "notify_status_update") {
			return p.parseNotifyStatusUpdate(message), nil
		}
		if strings.Contains(kpStatusString, "result") {
			return p.parseResult(message), nil
		}
	}
	return nil, nil
}

func (p *statePublisher) parseNotifyStatusUpdate(message []byte) []*models.PrinterStatus {
	var kpStatusUpdate *statusUpdate
	err := json.Unmarshal(message, &kpStatusUpdate)
	if err != nil {
		log.Println(err)
		return nil
	}

	status := make([]*models.PrinterStatus, 0)
	for _, p := range kpStatusUpdate.Params {
		if param, ok := p.(map[string]any); ok {
			for k, v := range param {
				status = append(status, &models.PrinterStatus{
					Name:  k,
					State: v,
				})
			}
		}
	}
	return status
}

func (p *statePublisher) parseResult(message []byte) []*models.PrinterStatus {
	var pkResult *result
	err := json.Unmarshal(message, &pkResult)
	if err != nil {
		log.Println(err)
		return nil
	}

	status := make([]*models.PrinterStatus, 0)
	for k, v := range pkResult.Result.Status {
		status = append(status, &models.PrinterStatus{
			Name:  k,
			State: v,
		})
	}
	return status
}

func (p *statePublisher) OnNewSub() {
	p.onNewSub <- struct{}{}
}

func (p *statePublisher) Close() {
	log.Println("state publisher Stop")
	p.conn.Close()
}
