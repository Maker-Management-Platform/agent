package printers

import (
	"context"
	"log"

	"github.com/duke-git/lancet/v2/maputil"
	"github.com/eduardooliveira/stLib/core/integrations/klipper"
	"github.com/eduardooliveira/stLib/core/models"
)

type Publisher interface {
	Start() error
	Out() <-chan []*models.PrinterStatus
	OnNewSub()
	Stop()
}

type stateManager struct {
	printer     *models.Printer
	publisher   Publisher
	subscribers []chan []*models.PrinterStatus
	sub         chan chan []*models.PrinterStatus
	unsub       chan (<-chan []*models.PrinterStatus)
}

func (s *stateManager) Subscribe(ctx context.Context) chan []*models.PrinterStatus {
	c := make(chan []*models.PrinterStatus, 10)
	s.sub <- c
	go func() {
		<-ctx.Done()
		s.Unsubscribe(c)
	}()
	return c
}
func (s *stateManager) Unsubscribe(channel chan []*models.PrinterStatus) {
	s.unsub <- channel
}

func (s *stateManager) run() {
	out := s.publisher.Out()
	for {
		select {
		case status := <-out:
			for _, c := range s.subscribers {
				c <- status
			}
		case channel := <-s.sub:
			s.subscribers = append(s.subscribers, channel)
			go s.publisher.OnNewSub()
			log.Println(s.printer.Name, " +Subs len: ", len(s.subscribers))
		case channel := <-s.unsub:
			for i, c := range s.subscribers {
				if c == channel {
					s.subscribers = append(s.subscribers[:i], s.subscribers[i+1:]...)
				}
			}
			log.Println(s.printer.Name, " -Subs len: ", len(s.subscribers))
			if len(s.subscribers) == 0 {
				s.publisher.Stop()
				stateManagers.Delete(s.printer.UUID)
				return
			}
		}
	}
}

var stateManagers = maputil.NewConcurrentMap[string, *stateManager](100)

func GetStateManager(p *models.Printer) (*stateManager, error) {
	if sm, ok := stateManagers.Get(p.UUID); ok {
		return sm, nil
	} else {
		sm := &stateManager{
			printer:     p,
			sub:         make(chan chan []*models.PrinterStatus),
			unsub:       make(chan (<-chan []*models.PrinterStatus)),
			subscribers: make([]chan []*models.PrinterStatus, 0),
		}
		stateManagers.Set(p.UUID, sm)
		if sm.publisher == nil {
			if sm.printer.Type == "klipper" {
				sm.publisher = klipper.GetStatePublisher(sm.printer)
				if err := sm.publisher.Start(); err != nil {
					return nil, err
				}
			}
		}
		go sm.run()
		return sm, nil
	}
}
