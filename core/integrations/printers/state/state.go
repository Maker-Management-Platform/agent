package state

import (
	"errors"
	"time"

	"github.com/eduardooliveira/stLib/core/entities"
	"github.com/eduardooliveira/stLib/core/integrations/models"
	"github.com/eduardooliveira/stLib/core/state"
)

var printerState = make(map[string]*models.PrinterState)

func init() {
	go runMock()
}

func runMock() {
	tick := 0
	time.Sleep(10 * time.Second)
	printerState["mock"] = &models.PrinterState{
		Printer: state.Printers["mock"],
		State:   "waiting_validation",
	}
	mock := printerState["mock"]
	for {
		tick += 1
		time.Sleep(1 * time.Second)
		if mock.State == "printing" {
			mock.Duration += 1
			mock.Progress += 10
		}
		if mock.Duration == 15 {
			mock.State = "waiting_validation"
			mock.Duration = 0
		}
		if tick == 10 {
			mock.State = "waiting_job"
		}
	}
}

func GetPrinterState(uuid string) (*models.PrinterState, bool) {
	p, ok := printerState[uuid]
	return p, ok
}

func PrintJob(uuid string, job *entities.PrintJob) error {
	ps, ok := printerState[uuid]
	if !ok {
		return errors.New("printer not found")
	}

	if ps.State != "waiting_job" {
		return errors.New("printer not awaiting job")
	}
	ps.Duration = 0
	ps.Progress = 0
	ps.JobUUID = job.UUID
	ps.State = "printing"

	return nil
}
