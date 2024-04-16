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
	printerState["mock"] = &models.PrinterState{
		Printer: state.Printers["mock"],
		State:   "awaiting_ok",
	}
	go runMock()
}

func runMock() {
	tick := 0
	mock := printerState["mock"]
	for {
		tick += 1
		time.Sleep(1 * time.Second)
		if mock.State == "printing" {
			mock.Duration += 1
		}
		if mock.Duration == 15 {
			mock.State = "awaiting_ok"
			mock.Duration = 0
		}
		if tick == 10 {
			mock.State = "awaiting_job"
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

	if ps.State != "awaiting_job" {
		return errors.New("printer not awaiting job")
	}

	ps.JobUUID = job.UUID
	ps.State = "printing"

	return nil
}
