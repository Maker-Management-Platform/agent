package system

import (
	"log"

	"github.com/eduardooliveira/stLib/core/events"
)

type systemEvent struct {
	Name  string `json:"name"`
	State any    `json:"state"`
}

type eventManagement struct {
}

func (em *eventManagement) Start() error {
	return nil
}
func (em *eventManagement) Stop() error {
	return nil
}
func (em *eventManagement) OnNewSub() error {
	return nil
}

func (em *eventManagement) Read() chan *events.Message {
	rtn := make(chan *events.Message, 1)
	eventName := "system.state"
	go func() {
		for {
			m := <-systemEvents
			select {
			case rtn <- &events.Message{
				Event: eventName,
				Data:  m,
			}:
				log.Println("event sent")
			default:
				log.Println("status update channel full")
			}
		}
	}()

	return rtn
}

var eventManager *eventManagement
var systemEvents chan *systemEvent

func GetEventPublisher() *eventManagement {
	return eventManager
}

func init() {
	eventManager = &eventManagement{}
	systemEvents = make(chan *systemEvent, 100)
}

func Publish(name string, data any) {
	select {
	case systemEvents <- &systemEvent{
		Name:  name,
		State: data,
	}:
	default:
		log.Println("dropped system event")
	}
}
