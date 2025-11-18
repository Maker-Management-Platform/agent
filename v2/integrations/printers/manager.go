package printers

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	v1octoprint "github.com/eduardooliveira/stLib/core/integrations/octorpint"
	v1klipper "github.com/eduardooliveira/stLib/core/integrations/klipper"
	v2events "github.com/eduardooliveira/stLib/v2/events"
)

// PrinterType represents the type of printer
type PrinterType string

const (
	PrinterTypeOctoPrint PrinterType = "octoprint"
	PrinterTypeKlipper   PrinterType = "klipper"
	PrinterTypeMoonraker PrinterType = "moonraker"
)

// PrinterState represents the state of a printer
type PrinterState string

const (
	StateOffline     PrinterState = "offline"
	StateIdle        PrinterState = "idle"
	StatePrinting    PrinterState = "printing"
	StatePaused      PrinterState = "paused"
	StateError       PrinterState = "error"
	StateConnecting  PrinterState = "connecting"
)

// Printer represents a 3D printer
type Printer struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        PrinterType            `json:"type"`
	URL         string                 `json:"url"`
	APIKey      string                 `json:"apiKey,omitempty"`
	State       PrinterState           `json:"state"`
	Temperature *TemperatureInfo       `json:"temperature,omitempty"`
	Progress    *PrintProgress         `json:"progress,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	LastUpdate  time.Time              `json:"lastUpdate"`
}

// TemperatureInfo holds temperature data
type TemperatureInfo struct {
	HotendTemp   float64 `json:"hotendTemp"`
	HotendTarget float64 `json:"hotendTarget"`
	BedTemp      float64 `json:"bedTemp"`
	BedTarget    float64 `json:"bedTarget"`
}

// PrintProgress holds print job progress
type PrintProgress struct {
	Completion      float64 `json:"completion"`
	PrintTime       int     `json:"printTime"`
	PrintTimeLeft   int     `json:"printTimeLeft"`
	CurrentFile     string  `json:"currentFile"`
}

// PrinterAdapter interface for different printer types
type PrinterAdapter interface {
	Connect(ctx context.Context) error
	Disconnect() error
	GetState(ctx context.Context) (*Printer, error)
	SendGCode(ctx context.Context, command string) error
	UploadFile(ctx context.Context, filename string, data []byte) error
	StartPrint(ctx context.Context, filename string) error
	PausePrint(ctx context.Context) error
	ResumePrint(ctx context.Context) error
	CancelPrint(ctx context.Context) error
}

// PrinterManager manages all printer connections
type PrinterManager struct {
	printers    map[string]*Printer
	adapters    map[string]PrinterAdapter
	mutex       sync.RWMutex
	eventMgr    *v2events.EventManager
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
}

// NewPrinterManager creates a new printer manager
func NewPrinterManager(eventMgr *v2events.EventManager) *PrinterManager {
	ctx, cancel := context.WithCancel(context.Background())

	return &PrinterManager{
		printers:    make(map[string]*Printer),
		adapters:    make(map[string]PrinterAdapter),
		eventMgr:    eventMgr,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Start begins monitoring all printers
func (m *PrinterManager) Start() {
	m.wg.Add(1)
	go m.monitorPrinters()
}

// Stop gracefully stops the printer manager
func (m *PrinterManager) Stop() {
	m.cancel()
	m.wg.Wait()

	// Disconnect all printers
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, adapter := range m.adapters {
		adapter.Disconnect()
	}
}

// AddPrinter adds a new printer to the manager
func (m *PrinterManager) AddPrinter(printer *Printer) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.printers[printer.ID]; exists {
		return fmt.Errorf("printer %s already exists", printer.ID)
	}

	// Create adapter based on printer type
	var adapter PrinterAdapter
	var err error

	switch printer.Type {
	case PrinterTypeOctoPrint:
		adapter, err = NewOctoPrintAdapter(printer)
	case PrinterTypeKlipper:
		adapter, err = NewKlipperAdapter(printer)
	default:
		return fmt.Errorf("unsupported printer type: %s", printer.Type)
	}

	if err != nil {
		return fmt.Errorf("failed to create adapter: %w", err)
	}

	m.printers[printer.ID] = printer
	m.adapters[printer.ID] = adapter

	slog.Info("Printer added", "id", printer.ID, "type", printer.Type)

	return nil
}

// RemovePrinter removes a printer from the manager
func (m *PrinterManager) RemovePrinter(printerID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	adapter, exists := m.adapters[printerID]
	if !exists {
		return fmt.Errorf("printer %s not found", printerID)
	}

	adapter.Disconnect()
	delete(m.printers, printerID)
	delete(m.adapters, printerID)

	slog.Info("Printer removed", "id", printerID)

	return nil
}

// GetPrinter returns a printer by ID
func (m *PrinterManager) GetPrinter(printerID string) (*Printer, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	printer, exists := m.printers[printerID]
	if !exists {
		return nil, fmt.Errorf("printer %s not found", printerID)
	}

	return printer, nil
}

// ListPrinters returns all printers
func (m *PrinterManager) ListPrinters() []*Printer {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	printers := make([]*Printer, 0, len(m.printers))
	for _, printer := range m.printers {
		printers = append(printers, printer)
	}

	return printers
}

// SendGCode sends a G-code command to a printer
func (m *PrinterManager) SendGCode(printerID string, command string) error {
	m.mutex.RLock()
	adapter, exists := m.adapters[printerID]
	m.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("printer %s not found", printerID)
	}

	ctx, cancel := context.WithTimeout(m.ctx, 30*time.Second)
	defer cancel()

	return adapter.SendGCode(ctx, command)
}

// UploadFile uploads a file to a printer
func (m *PrinterManager) UploadFile(printerID string, filename string, data []byte) error {
	m.mutex.RLock()
	adapter, exists := m.adapters[printerID]
	m.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("printer %s not found", printerID)
	}

	ctx, cancel := context.WithTimeout(m.ctx, 5*time.Minute)
	defer cancel()

	return adapter.UploadFile(ctx, filename, data)
}

// StartPrint starts a print job
func (m *PrinterManager) StartPrint(printerID string, filename string) error {
	m.mutex.RLock()
	adapter, exists := m.adapters[printerID]
	m.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("printer %s not found", printerID)
	}

	ctx, cancel := context.WithTimeout(m.ctx, 30*time.Second)
	defer cancel()

	return adapter.StartPrint(ctx, filename)
}

// monitorPrinters periodically updates printer states
func (m *PrinterManager) monitorPrinters() {
	defer m.wg.Done()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.updateAllPrinters()
		}
	}
}

// updateAllPrinters updates the state of all printers
func (m *PrinterManager) updateAllPrinters() {
	m.mutex.RLock()
	adapters := make(map[string]PrinterAdapter, len(m.adapters))
	for id, adapter := range m.adapters {
		adapters[id] = adapter
	}
	m.mutex.RUnlock()

	for printerID, adapter := range adapters {
		go m.updatePrinter(printerID, adapter)
	}
}

// updatePrinter updates a single printer's state
func (m *PrinterManager) updatePrinter(printerID string, adapter PrinterAdapter) {
	ctx, cancel := context.WithTimeout(m.ctx, 10*time.Second)
	defer cancel()

	printer, err := adapter.GetState(ctx)
	if err != nil {
		slog.Debug("Failed to get printer state", "printer", printerID, "error", err)

		// Update printer to offline state
		m.mutex.Lock()
		if p, exists := m.printers[printerID]; exists {
			p.State = StateOffline
			p.LastUpdate = time.Now()
		}
		m.mutex.Unlock()
		return
	}

	printer.LastUpdate = time.Now()

	// Update stored printer state
	m.mutex.Lock()
	m.printers[printerID] = printer
	m.mutex.Unlock()

	// Broadcast printer status update
	if m.eventMgr != nil {
		m.eventMgr.Broadcast(v2events.EventPrinterStatus, map[string]interface{}{
			"id":      printerID,
			"printer": printer,
		})
	}
}

// LoadV1Printers loads printers from v1 integration layer
func (m *PrinterManager) LoadV1Printers() error {
	// This would load printers from v1's state
	// For now, we'll create a placeholder

	slog.Info("Loading v1 printers (placeholder)")
	return nil
}

// AdaptV1OctoPrint creates a v2 printer from v1 OctoPrint
func AdaptV1OctoPrint(v1Printer *v1octoprint.OctoPrint) *Printer {
	state := StateOffline
	if v1Printer != nil {
		// Map v1 state to v2 state
		// This is a simplified mapping
		state = StateIdle
	}

	return &Printer{
		ID:         v1Printer.Name,
		Name:       v1Printer.Name,
		Type:       PrinterTypeOctoPrint,
		URL:        v1Printer.URL,
		APIKey:     v1Printer.APIKey,
		State:      state,
		LastUpdate: time.Now(),
	}
}

// AdaptV1Klipper creates a v2 printer from v1 Klipper
func AdaptV1Klipper(v1Printer *v1klipper.Klipper) *Printer {
	state := StateOffline
	if v1Printer != nil {
		state = StateIdle
	}

	return &Printer{
		ID:         v1Printer.Name,
		Name:       v1Printer.Name,
		Type:       PrinterTypeKlipper,
		URL:        v1Printer.URL,
		State:      state,
		LastUpdate: time.Now(),
	}
}
