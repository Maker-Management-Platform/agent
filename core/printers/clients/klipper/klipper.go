package klipper

import (
	coreEntities "github.com/eduardooliveira/stLib/core/entities"
	"github.com/eduardooliveira/stLib/core/events"
	printerEntities "github.com/eduardooliveira/stLib/core/printers/entities"
)

func NewPrinter(config *printerEntities.Config) printerEntities.Printer {
	return KlipperPrinter{
		config:                config,
		state:                 &printerEntities.State{},
		bedState:              &printerEntities.TemperatureStatus{},
		hotEndState:           []*printerEntities.TemperatureStatus{&printerEntities.TemperatureStatus{}},
		jobState:              &printerEntities.JobStatus{},
		bedChangeListeners:    make(map[chan *printerEntities.TemperatureStatus]struct{}, 0),
		hotEndChangeListeners: make(map[chan []*printerEntities.TemperatureStatus]struct{}, 0),
		jobChangeListeners:    make(map[chan *printerEntities.JobStatus]struct{}, 0),
		deltaListeners:        make(map[chan []*events.Message]struct{}, 0),
	}
}

func (kp KlipperPrinter) GetConfig() *printerEntities.Config {
	return kp.config
}

func (kp KlipperPrinter) UploadFile(asset *coreEntities.ProjectAsset) error {
	return kp.serverFilesUpload(asset)
}

func (kp KlipperPrinter) TestConnection() error {

	r, err := kp.serverInfo()

	if err != nil {
		kp.state.ConnectionStatus = "disconnected"
		return err
	}

	kp.state.Version = r.APIVersionString
	kp.state.ConnectionStatus = "connected"

	return nil
}

func (kp KlipperPrinter) GetState() *printerEntities.State {
	return kp.state
}

func (kp KlipperPrinter) OnBedChange(out chan *printerEntities.TemperatureStatus) {
	kp.bedChangeListeners[out] = struct{}{}
}

func (kp KlipperPrinter) OnHotEndChange(out chan []*printerEntities.TemperatureStatus) {
	kp.hotEndChangeListeners[out] = struct{}{}
}

func (kp KlipperPrinter) OnJobChange(out chan *printerEntities.JobStatus) {
	kp.jobChangeListeners[out] = struct{}{}
}

func (kp KlipperPrinter) OnDeltaChange(out chan []*events.Message) {
	kp.deltaListeners[out] = struct{}{}
}

func broadcast[T any](outs map[chan T]struct{}, data T) {
	for out := range outs {
		select {
		case out <- data:
		default:
			delete(outs, out)
		}
	}
}
