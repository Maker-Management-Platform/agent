package state

import (
	"log"
	"strings"
	"time"

	"github.com/eduardooliveira/stLib/core/events"
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
				jsu := &entities.JobStatusUpdate{
					UUID: p.UUID,
				}
				var sd any
				var ok bool
				var err error
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
