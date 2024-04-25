package state

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	coreEntities "github.com/eduardooliveira/stLib/core/entities"
	"github.com/eduardooliveira/stLib/core/events"
	"github.com/eduardooliveira/stLib/core/integrations/models"
	printersState "github.com/eduardooliveira/stLib/core/integrations/printers/state"
	"github.com/eduardooliveira/stLib/core/printqueue/entities"
)

var updateInterval = 10 * time.Second
var triggerUpdateChannel = make(chan any)
var printerQueue = make(map[string]time.Time)

func init() {
	go runUpdate()
}

func triggerUpdate() {
	go func() {
		triggerUpdateChannel <- struct{}{}
	}()
}

func asyncUpdate() {
	ticker := time.NewTicker(updateInterval)
	for {
		select {
		case <-ticker.C:
		case <-triggerUpdateChannel:
		}
		updatePrintQueue()
	}
}

func updatePrintQueue() {
	if len(printQueue) == 0 {
		return
	}
	updateEvents := make([]*events.Message, 0)
	for _, job := range printQueue {

		if job.State == "canceled" || job.State == "failed" {
			continue
		}
		jsu := &entities.JobStatusUpdate{
			UUID: job.UUID,
		}

		printerState, err := matchPrintStateToJob(job)

		if err != nil {
			log.Println("printer not found ", job.PrinterUUID)
			jsu.Error = err.Error()
			updateEvents = append(updateEvents, &events.Message{
				Event: job.UUID,
				Data:  jsu,
			})
			continue
		}

		switch printerState.State {
		case "printing":
			if printerState.JobUUID == job.UUID {

			}
		}

	}
}

func runUpdate() {
	for {
		if len(printQueue) > 0 {

			startDelta := time.Now()
			updateEvents := make([]*events.Message, 0)
			for _, p := range printQueue {

				if p.State != "queued" && p.State != "printing" {
					continue
				}

				jsu := &entities.JobStatusUpdate{
					UUID: p.UUID,
				}
				pState, err := matchPrintStateToJob(p)

				if err != nil {
					log.Println("printer not found ", p.PrinterUUID)
					jsu.Error = err.Error()
					updateEvents = append(updateEvents, &events.Message{
						Event: p.UUID,
						Data:  jsu,
					})
					continue
				}

				switch pState.State {
				case "waiting_job":
					err = printersState.PrintJob(pState.Printer.UUID, p)
					if err != nil {
						log.Println("error printing job ", p.UUID, err)
						jsu.Error = err.Error()
						updateEvents = append(updateEvents, &events.Message{
							Event: p.UUID,
							Data:  jsu,
						})
						continue
					}
					startDelta = time.Now()
					start(p.UUID)

				case "waiting_validation":
					if pState.JobUUID == p.UUID {
						finish(p.UUID)
					}
				}

				jsu.StartAt = startDelta
				jsu.EndAt = jsu.StartAt.Add(jsu.Duration)

				startDelta = jsu.EndAt
				updateEvents = append(updateEvents, &events.Message{
					Event: p.UUID,
					Data:  jsu,
				})
			}
			Publish("job.update", updateEvents, true)
		}

		time.Sleep(updateInterval)
	}
}

func matchPrintStateToJob(p *coreEntities.PrintJob) (*models.PrinterState, error) {
	if p.PrinterUUID != "" {
		pState, ok := printersState.GetPrinterState(p.PrinterUUID)
		if !ok {
			return nil, errors.New("printer not found")
		}
		return pState, nil
	}
	// handle tags
	return nil, errors.New("no match found")
}

func sliceDuration(s *coreEntities.ProjectAsset) (time.Duration, error) {
	var sd any
	var ok bool
	if sd, ok = s.Properties["estimated printing time (normal mode)"]; !ok {
		return 0, fmt.Errorf("estimated printing time (normal mode) not defined %s", s.Name)
	}
	var sds string
	if sds, ok = sd.(string); !ok {
		return 0, fmt.Errorf("estimated printing time (normal mode) not string %s", s.Name)
	}
	sds = strings.ReplaceAll(sds, " ", "")
	var rtn time.Duration
	var err error
	if rtn, err = time.ParseDuration(sds); err != nil {
		return 0, errors.Join(err, fmt.Errorf("estimated printing time (normal mode) not valid %s", s.Name))
	}
	return rtn, nil
}
