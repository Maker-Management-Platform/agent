package state

import (
	"log"

	"github.com/eduardooliveira/stLib/core/events"
)

type printQueueEvent struct {
	Name string `json:"name"`
	Data any    `json:"data"`
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
	queue, err := GetPrintQueue(nil)
	if err != nil {
		log.Println(err)
	}
	eventQueue <- &printQueueEvent{
		Name: "queue.update",
		Data: queue,
	}
	return nil
}

func (em *eventManagement) Read() chan *events.Message {
	rtn := make(chan *events.Message, 1)
	eventName := "printQueue"
	go func() {
		for {
			m := <-eventQueue
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
var eventQueue chan *printQueueEvent

func GetEventPublisher() *eventManagement {
	return eventManager
}

func init() {
	eventManager = &eventManagement{}
	eventQueue = make(chan *printQueueEvent, 100)
}

func Publish(name string, data any) {
	select {
	case eventQueue <- &printQueueEvent{
		Name: name,
		Data: data,
	}:
	default:
		log.Println("dropped print queue event")
	}
}
