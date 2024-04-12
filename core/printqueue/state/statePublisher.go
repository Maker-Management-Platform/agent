package state

import (
	"fmt"
	"log"

	"github.com/eduardooliveira/stLib/core/events"
)

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
	Publish("queue.update", queue, false)
	return nil
}

func (em *eventManagement) Read() chan *events.Message {
	rtn := make(chan *events.Message, 1)
	eventName := "printQueue.%s"
	go func() {
		for {
			m := <-eventQueue
			select {
			case rtn <- &events.Message{
				Event:  fmt.Sprintf(eventName, m.Event),
				Data:   m.Data,
				Unpack: m.Unpack,
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
var eventQueue chan *events.Message

func GetEventPublisher() *eventManagement {
	return eventManager
}

func init() {
	eventManager = &eventManagement{}
	eventQueue = make(chan *events.Message, 100)
}

func Publish(name string, data any, unpack bool) {
	select {
	case eventQueue <- &events.Message{
		Event:  name,
		Data:   data,
		Unpack: unpack,
	}:
	default:
		log.Println("dropped print queue event")
	}
}
