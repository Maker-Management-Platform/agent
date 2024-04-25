package state

import (
	"slices"

	"github.com/eduardooliveira/stLib/core/data/database"
	"github.com/eduardooliveira/stLib/core/entities"
)

var printQueue = make([]*entities.PrintJob, 0)

func Enqueue(job *entities.PrintJob) error {
	err := database.InsertPrintJob(job)
	if err != nil {
		return err
	}
	printQueue = append(printQueue, job)
	Publish("queue.update", printQueue, false)
	return nil
}

func Cancel(jobUUID string) ([]*entities.PrintJob, error) {
	return changePrintJobState(jobUUID, "cancelled")
}

func Validate(jobUUID string, state string) ([]*entities.PrintJob, error) {
	return changePrintJobState(jobUUID, "validated")
}

func start(jobUUID string) ([]*entities.PrintJob, error) {
	return changePrintJobState(jobUUID, "printing")
}

func finish(jobUUID string) ([]*entities.PrintJob, error) {
	return changePrintJobState(jobUUID, "finish")
}

func success(jobUUID string) ([]*entities.PrintJob, error) {
	return changePrintJobState(jobUUID, "success")
}

func fail(jobUUID string) ([]*entities.PrintJob, error) {
	return changePrintJobState(jobUUID, "failed")
}

func changePrintJobState(jobUUID string, state string) ([]*entities.PrintJob, error) {
	err := database.ChangePrintJobState(jobUUID, state)
	if err != nil {
		return nil, err
	}
	printQueue, err = database.GetPrintJobs()
	Publish("queue.update", printQueue, false)
	return printQueue, err
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
