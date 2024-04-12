package state

import (
	"slices"

	"github.com/eduardooliveira/stLib/core/data/database"
	"github.com/eduardooliveira/stLib/core/entities"
)

var printQueue = make([]*entities.PrintJob, 0)

func AddPrintJob(job *entities.PrintJob) error {
	err := database.InsertPrintJob(job)
	if err != nil {
		return err
	}
	printQueue = append(printQueue, job)
	Publish("queue.update", printQueue, false)
	return nil
}

func GetPrintQueue(states []string) ([]*entities.PrintJob, error) {
	if len(printQueue) == 0 {
		var err error
		printQueue, err = database.GetPrintJobs()
		if err != nil {
			return nil, err
		}
	}
	if len(states) != 0 {
		filtered := make([]*entities.PrintJob, 0)
		for _, job := range printQueue {
			if slices.Contains(states, job.State) {
				filtered = append(filtered, job)
			}
		}
		return filtered, nil
	}
	return printQueue, nil
}

func MovePrintJob(jobUUID string, pos int) ([]*entities.PrintJob, error) {
	var err error
	err = database.Move(jobUUID, pos)
	if err != nil {
		return nil, err
	}
	printQueue, err = database.GetPrintJobs()
	Publish("queue.update", printQueue, false)
	return printQueue, err
}
