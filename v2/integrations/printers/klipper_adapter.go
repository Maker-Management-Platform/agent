package printers

import (
	"context"
	"fmt"

	v1klipper "github.com/eduardooliveira/stLib/core/integrations/klipper"
)

// KlipperAdapter wraps v1 Klipper integration for v2
type KlipperAdapter struct {
	printer    *Printer
	v1Client   *v1klipper.Klipper
}

// NewKlipperAdapter creates a new Klipper adapter
func NewKlipperAdapter(printer *Printer) (*KlipperAdapter, error) {
	// Create v1 Klipper client
	v1Client := &v1klipper.Klipper{
		Name: printer.Name,
		URL:  printer.URL,
	}

	return &KlipperAdapter{
		printer:  printer,
		v1Client: v1Client,
	}, nil
}

// Connect establishes connection to the printer
func (a *KlipperAdapter) Connect(ctx context.Context) error {
	// Klipper/Moonraker connections are stateless (HTTP-based)
	// Just verify the connection works
	_, err := a.GetState(ctx)
	return err
}

// Disconnect closes the connection
func (a *KlipperAdapter) Disconnect() error {
	// Klipper is stateless, nothing to disconnect
	return nil
}

// GetState retrieves the current printer state
func (a *KlipperAdapter) GetState(ctx context.Context) (*Printer, error) {
	// Get state from v1 client
	state, err := a.v1Client.GetState()
	if err != nil {
		return nil, fmt.Errorf("failed to get Klipper state: %w", err)
	}

	// Map v1 state to v2 printer
	printer := &Printer{
		ID:     a.printer.ID,
		Name:   a.printer.Name,
		Type:   PrinterTypeKlipper,
		URL:    a.printer.URL,
		State:  mapKlipperState(state.State),
	}

	// Map temperature if available
	if state.Extruder != nil {
		printer.Temperature = &TemperatureInfo{
			HotendTemp:   state.Extruder.Temperature,
			HotendTarget: state.Extruder.Target,
		}
	}

	if state.Heater_bed != nil {
		if printer.Temperature == nil {
			printer.Temperature = &TemperatureInfo{}
		}
		printer.Temperature.BedTemp = state.Heater_bed.Temperature
		printer.Temperature.BedTarget = state.Heater_bed.Target
	}

	// Map progress if printing
	if state.DisplayStatus != nil {
		printer.Progress = &PrintProgress{
			Completion: state.DisplayStatus.Progress * 100,
		}

		if state.PrintStats != nil {
			printer.Progress.PrintTime = int(state.PrintStats.PrintDuration)
			printer.Progress.CurrentFile = state.PrintStats.Filename
		}
	}

	return printer, nil
}

// SendGCode sends a G-code command to the printer
func (a *KlipperAdapter) SendGCode(ctx context.Context, command string) error {
	return a.v1Client.SendGCode(command)
}

// UploadFile uploads a file to Klipper/Moonraker
func (a *KlipperAdapter) UploadFile(ctx context.Context, filename string, data []byte) error {
	return a.v1Client.UploadFile(filename, data)
}

// StartPrint starts printing a file
func (a *KlipperAdapter) StartPrint(ctx context.Context, filename string) error {
	return a.v1Client.StartPrint(filename)
}

// PausePrint pauses the current print
func (a *KlipperAdapter) PausePrint(ctx context.Context) error {
	return a.v1Client.PausePrint()
}

// ResumePrint resumes a paused print
func (a *KlipperAdapter) ResumePrint(ctx context.Context) error {
	return a.v1Client.ResumePrint()
}

// CancelPrint cancels the current print
func (a *KlipperAdapter) CancelPrint(ctx context.Context) error {
	return a.v1Client.CancelPrint()
}

// mapKlipperState maps Klipper state strings to v2 PrinterState
func mapKlipperState(stateText string) PrinterState {
	switch stateText {
	case "ready":
		return StateIdle
	case "printing":
		return StatePrinting
	case "paused":
		return StatePaused
	case "error":
		return StateError
	case "shutdown":
		return StateError
	case "startup":
		return StateConnecting
	default:
		return StateOffline
	}
}
