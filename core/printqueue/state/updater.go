package state

import (
	"errors"
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

func init() {
	go runUpdate()
}

func runUpdate() {
	for {
		if len(printQueue) > 0 {

			startDelta := time.Now()
			updateEvents := make([]*events.Message, 0)
			for _, p := range printQueue {

				if p.State != "queued" {
					continue
				}

				jsu := &entities.JobStatusUpdate{
					UUID: p.UUID,
				}
				pState, err := matchPrinterToJob(p)

				if err != nil {
					log.Println("printer not found ", p.PrinterUUID)
					jsu.Error = err.Error()
					updateEvents = append(updateEvents, &events.Message{
						Event: p.UUID,
						Data:  jsu,
					})
					continue
				}

				if pState.State == "awaiting_job" {
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
				}

				var sd any
				var ok bool
				if sd, ok = p.Slice.Properties["estimated printing time (normal mode)"]; !ok {
					log.Println("estimated printing time (normal mode) not defined ", p.Slice.Name)
					continue
				}
				var sds string
				if sds, ok = sd.(string); !ok {
					log.Println("estimated printing time (normal mode) not string ", p.Slice.Name)
					continue
				}
				sds = strings.ReplaceAll(sds, " ", "")
				if jsu.Duration, err = time.ParseDuration(sds); err != nil {
					log.Println("estimated printing time (normal mode) not valid ", p.Slice.Name, sds)
					continue
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

func matchPrinterToJob(p *coreEntities.PrintJob) (*models.PrinterState, error) {
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
