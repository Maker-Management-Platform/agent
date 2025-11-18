package printers

import (
	"context"
	"fmt"

	v1octoprint "github.com/eduardooliveira/stLib/core/integrations/octorpint"
)

// OctoPrintAdapter wraps v1 OctoPrint integration for v2
type OctoPrintAdapter struct {
	printer    *Printer
	v1Client   *v1octoprint.OctoPrint
}

// NewOctoPrintAdapter creates a new OctoPrint adapter
func NewOctoPrintAdapter(printer *Printer) (*OctoPrintAdapter, error) {
	// Create v1 OctoPrint client
	v1Client := &v1octoprint.OctoPrint{
		Name:   printer.Name,
		URL:    printer.URL,
		APIKey: printer.APIKey,
	}

	return &OctoPrintAdapter{
		printer:  printer,
		v1Client: v1Client,
	}, nil
}

// Connect establishes connection to the printer
func (a *OctoPrintAdapter) Connect(ctx context.Context) error {
	// OctoPrint connections are stateless (HTTP-based)
	// Just verify the connection works
	_, err := a.GetState(ctx)
	return err
}

// Disconnect closes the connection
func (a *OctoPrintAdapter) Disconnect() error {
	// OctoPrint is stateless, nothing to disconnect
	return nil
}

// GetState retrieves the current printer state
func (a *OctoPrintAdapter) GetState(ctx context.Context) (*Printer, error) {
	// Get state from v1 client
	state, err := a.v1Client.GetState()
	if err != nil {
		return nil, fmt.Errorf("failed to get OctoPrint state: %w", err)
	}

	// Map v1 state to v2 printer
	printer := &Printer{
		ID:     a.printer.ID,
		Name:   a.printer.Name,
		Type:   PrinterTypeOctoPrint,
		URL:    a.printer.URL,
		APIKey: a.printer.APIKey,
		State:  mapOctoPrintState(state.State.Text),
	}

	// Map temperature if available
	if state.Temperature != nil && len(state.Temperature) > 0 {
		printer.Temperature = &TemperatureInfo{}

		// Get tool temperature (hotend)
		if tool, ok := state.Temperature["tool0"]; ok {
			if actual, ok := tool["actual"].(float64); ok {
				printer.Temperature.HotendTemp = actual
			}
			if target, ok := tool["target"].(float64); ok {
				printer.Temperature.HotendTarget = target
			}
		}

		// Get bed temperature
		if bed, ok := state.Temperature["bed"]; ok {
			if actual, ok := bed["actual"].(float64); ok {
				printer.Temperature.BedTemp = actual
			}
			if target, ok := bed["target"].(float64); ok {
				printer.Temperature.BedTarget = target
			}
		}
	}

	// Map progress if printing
	if state.Progress != nil {
		printer.Progress = &PrintProgress{
			Completion:    state.Progress.Completion,
			PrintTime:     state.Progress.PrintTime,
			PrintTimeLeft: state.Progress.PrintTimeLeft,
		}

		if state.Job != nil && state.Job.File != nil {
			printer.Progress.CurrentFile = state.Job.File.Name
		}
	}

	return printer, nil
}

// SendGCode sends a G-code command to the printer
func (a *OctoPrintAdapter) SendGCode(ctx context.Context, command string) error {
	return a.v1Client.SendGCode(command)
}

// UploadFile uploads a file to OctoPrint
func (a *OctoPrintAdapter) UploadFile(ctx context.Context, filename string, data []byte) error {
	return a.v1Client.UploadFile(filename, data)
}

// StartPrint starts printing a file
func (a *OctoPrintAdapter) StartPrint(ctx context.Context, filename string) error {
	return a.v1Client.StartPrint(filename)
}

// PausePrint pauses the current print
func (a *OctoPrintAdapter) PausePrint(ctx context.Context) error {
	return a.v1Client.PausePrint()
}

// ResumePrint resumes a paused print
func (a *OctoPrintAdapter) ResumePrint(ctx context.Context) error {
	return a.v1Client.ResumePrint()
}

// CancelPrint cancels the current print
func (a *OctoPrintAdapter) CancelPrint(ctx context.Context) error {
	return a.v1Client.CancelPrint()
}

// mapOctoPrintState maps OctoPrint state strings to v2 PrinterState
func mapOctoPrintState(stateText string) PrinterState {
	switch stateText {
	case "Operational":
		return StateIdle
	case "Printing":
		return StatePrinting
	case "Paused":
		return StatePaused
	case "Error":
		return StateError
	case "Offline":
		return StateOffline
	case "Connecting":
		return StateConnecting
	default:
		return StateOffline
	}
}
