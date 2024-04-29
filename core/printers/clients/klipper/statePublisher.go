package klipper

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/duke-git/lancet/v2/maputil"
	models "github.com/eduardooliveira/stLib/core/entities"
	"github.com/eduardooliveira/stLib/core/events"
	printerModels "github.com/eduardooliveira/stLib/core/integrations/models"
	printerEntities "github.com/eduardooliveira/stLib/core/printers/entities"
	"github.com/gorilla/websocket"
)

type statePublisher struct {
	printer  *KlipperPrinter
	out      chan *events.Message
	onNewSub chan struct{}
	done     chan struct{}
	conn     *websocket.Conn
}

func GetStatePublisher(printer *models.Printer) *statePublisher {
	kp := &KlipperPrinter{printer}
	return &statePublisher{
		printer:  kp,
		out:      make(chan *events.Message),
		onNewSub: make(chan struct{}),
		done:     make(chan struct{}),
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
func (p *statePublisher) Read() chan *events.Message {
	rtn := make(chan *events.Message, 10)
	eventName := fmt.Sprintf("printer.update.%s", p.printer.config.UUID)
	go func() {
		for {
			select {
			case <-p.done:
				return
			case <-p.onNewSub:
				p.conn.WriteMessage(websocket.TextMessage, []byte("{\"jsonrpc\": \"2.0\",\"method\": \"printer.objects.query\",\"params\": {\"objects\": {\"extruder\": null,\"heater_bed\": null, \"print_stats\":null, \"display_status\": null}},\"id\": 2}"))
			default:
				_, message, err := p.conn.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					//TODO: implement reconnect
					close(rtn)
					return
				}

				kpStatusString := string(message)

				if strings.Contains(kpStatusString, "notify_proc_stat_update") {
					continue
				}

				//log.Println(p.printer.Name, "status update:", kpStatusString)

				if strings.Contains(kpStatusString, "notify_status_update") {
					//log.Println(p.printer.Name, "status update:", kpStatusString)
					select {
					case rtn <- &events.Message{
						Event:  eventName,
						Data:   p.parseNotifyStatusUpdate(message),
						Unpack: true,
					}:
					default:
						log.Println("status update channel full")
					}
				}
				if strings.Contains(kpStatusString, "result") {
					log.Println(p.printer.config.Name, "status update:", kpStatusString)
					select {
					case rtn <- &events.Message{
						Event:  eventName,
						Data:   p.parseResult(message),
						Unpack: true,
					}:
					default:
						log.Println("status update channel full")
					}
				}
			}
		}
	}()

	return rtn
}

func (p *statePublisher) OnNewSub() error {
	p.onNewSub <- struct{}{}
	return nil
}

func (p *statePublisher) Stop() error {
	log.Println("state publisher Stop")
	p.done <- struct{}{}
	p.conn.Close()
	return nil
}

func (kp KlipperPrinter) parseStatus(name string, state map[string]any, status map[string]*events.Message) {

	switch name {
	case "heater_bed":
		status["bed"] = &events.Message{
			Event: "bed",
			Data:  &printerModels.TemperatureStatus{},
		}
		handleThermalValue("bed", state, status, kp.bedState)
		broadcast[*printerEntities.TemperatureStatus](kp.bedChangeListeners, kp.bedState)
	case "extruder":
		status["extruder"] = &events.Message{
			Event: "extruder",
			Data:  &printerModels.TemperatureStatus{},
		}
		handleThermalValue("extruder", state, status, kp.bedState)
		broadcast[[]*printerEntities.TemperatureStatus](kp.hotEndChangeListeners, kp.hotEndState)
	case "print_stats":
		var ok bool
		_, ok = status["job_status"]
		if !ok {
			status["job_status"] = &events.Message{
				Event: "job_status",
				Data:  &printerModels.JobStatus{},
			}
		}
		current := status["job_status"].Data.(*printerModels.JobStatus)
		if v, ok := state["total_duration"].(float64); ok {
			current.TotalDuration = v
			kp.jobState.TotalDuration = v
		}
		if v, ok := state["filename"].(string); ok {
			current.FileName = v
			kp.jobState.FileName = v
		}
		broadcast[*printerEntities.JobStatus](kp.jobChangeListeners, kp.jobState)
	case "display_status":
		var ok bool
		_, ok = status["job_status"]
		if !ok {
			status["job_status"] = &events.Message{
				Event: "job_status",
				Data:  &printerModels.JobStatus{},
			}
		}
		current := status["job_status"].Data.(*printerModels.JobStatus)
		if v, ok := state["message"].(string); ok {
			current.Message = v
			kp.jobState.Message = v
		}
		if v, ok := state["progress"].(float64); ok {
			current.Progress = v
			kp.jobState.Progress = v
		}
		broadcast[*printerEntities.JobStatus](kp.jobChangeListeners, kp.jobState)
	}

}

func handleThermalValue(key string, values map[string]any, status map[string]*events.Message, temperatureStatus *printerEntities.TemperatureStatus) {
	if v, ok := values["temperature"].(float64); ok {
		status[key].Data.(*printerModels.TemperatureStatus).Temperature = v
		temperatureStatus.Temperature = v
	}
	if v, ok := values["target"].(float64); ok {
		status[key].Data.(*printerModels.TemperatureStatus).Target = v
		temperatureStatus.Target = v
	}
	if v, ok := values["power"].(float64); ok {
		status[key].Data.(*printerModels.TemperatureStatus).Power = v
		temperatureStatus.Power = v
	}
}

func (p *statePublisher) parseNotifyStatusUpdate(message []byte) []*events.Message {
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
				parseStatus(k, v.(map[string]any), status)
			}
		}
	}
	return maputil.Values(status)
}

func (p *statePublisher) parseResult(message []byte) []*events.Message {
	var pkResult *result
	err := json.Unmarshal(message, &pkResult)
	if err != nil {
		log.Println(err)
		return nil
	}

	status := make(map[string]*events.Message, 0)
	for k, v := range pkResult.Result.Status {
		parseStatus(k, v, status)
	}
	return maputil.Values(status)
}
