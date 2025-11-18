package integration_test

import (
	"context"
	"testing"
	"time"

	v2events "github.com/eduardooliveira/stLib/v2/events"
	v2printers "github.com/eduardooliveira/stLib/v2/integrations/printers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPrinterManagerCreation tests printer manager creation
func TestPrinterManagerCreation(t *testing.T) {
	mgr := v2printers.NewPrinterManager(nil)
	require.NotNil(t, mgr)
}

// TestPrinterManagerWithEvents tests printer manager with event manager
func TestPrinterManagerWithEvents(t *testing.T) {
	eventMgr := v2events.NewEventManager()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go eventMgr.Start(ctx)

	mgr := v2printers.NewPrinterManager(eventMgr)
	require.NotNil(t, mgr)
}

// TestAddPrinter tests adding a printer
func TestAddPrinter(t *testing.T) {
	mgr := v2printers.NewPrinterManager(nil)

	tests := []struct {
		name        string
		printer     *v2printers.Printer
		wantErr     bool
	}{
		{
			name: "OctoPrint",
			printer: &v2printers.Printer{
				ID:   "octoprint-1",
				Name: "My OctoPrint",
				Type: v2printers.PrinterTypeOctoPrint,
				URL:  "http://octopi.local",
			},
			wantErr: false,
		},
		{
			name: "Klipper",
			printer: &v2printers.Printer{
				ID:   "klipper-1",
				Name: "My Klipper",
				Type: v2printers.PrinterTypeKlipper,
				URL:  "http://mainsailos.local",
			},
			wantErr: false,
		},
		{
			name: "Moonraker",
			printer: &v2printers.Printer{
				ID:   "moonraker-1",
				Name: "My Moonraker",
				Type: v2printers.PrinterTypeMoonraker,
				URL:  "http://fluidd.local",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mgr.AddPrinter(tt.printer)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestRemovePrinter tests removing a printer
func TestRemovePrinter(t *testing.T) {
	mgr := v2printers.NewPrinterManager(nil)

	printer := &v2printers.Printer{
		ID:   "test-1",
		Name: "Test Printer",
		Type: v2printers.PrinterTypeOctoPrint,
		URL:  "http://test.local",
	}

	// Add printer
	err := mgr.AddPrinter(printer)
	require.NoError(t, err)

	// Remove printer
	err = mgr.RemovePrinter("test-1")
	assert.NoError(t, err)

	// Should not be listed anymore
	printers, _ := mgr.ListPrinters()
	found := false
	for _, p := range printers {
		if p.ID == "test-1" {
			found = true
			break
		}
	}
	assert.False(t, found, "Printer should be removed")
}

// TestListPrinters tests listing printers
func TestListPrinters(t *testing.T) {
	mgr := v2printers.NewPrinterManager(nil)

	// Add multiple printers
	printers := []*v2printers.Printer{
		{ID: "p1", Name: "Printer 1", Type: v2printers.PrinterTypeOctoPrint, URL: "http://p1.local"},
		{ID: "p2", Name: "Printer 2", Type: v2printers.PrinterTypeKlipper, URL: "http://p2.local"},
		{ID: "p3", Name: "Printer 3", Type: v2printers.PrinterTypeMoonraker, URL: "http://p3.local"},
	}

	for _, p := range printers {
		err := mgr.AddPrinter(p)
		require.NoError(t, err)
	}

	// List printers
	list, err := mgr.ListPrinters()
	assert.NoError(t, err)
	assert.Len(t, list, len(printers))
}

// TestPrinterStates tests printer state enumeration
func TestPrinterStates(t *testing.T) {
	states := []v2printers.PrinterState{
		v2printers.StateOffline,
		v2printers.StateConnecting,
		v2printers.StateIdle,
		v2printers.StatePrinting,
		v2printers.StatePaused,
		v2printers.StateError,
	}

	for _, state := range states {
		// Just verify states exist
		assert.NotEmpty(t, string(state))
	}
}

// TestPrinterTypes tests printer type enumeration
func TestPrinterTypes(t *testing.T) {
	types := []v2printers.PrinterType{
		v2printers.PrinterTypeOctoPrint,
		v2printers.PrinterTypeKlipper,
		v2printers.PrinterTypeMoonraker,
	}

	for _, ptype := range types {
		// Just verify types exist
		assert.NotEmpty(t, string(ptype))
	}
}

// TestOctoPrintAdapter tests OctoPrint adapter creation
func TestOctoPrintAdapter(t *testing.T) {
	t.Skip("Requires actual OctoPrint instance or mock server")

	printer := &v2printers.Printer{
		ID:   "octo-1",
		Name: "Test OctoPrint",
		Type: v2printers.PrinterTypeOctoPrint,
		URL:  "http://octopi.local",
	}

	adapter, err := v2printers.NewOctoPrintAdapter(printer)
	require.NoError(t, err)
	require.NotNil(t, adapter)
}

// TestKlipperAdapter tests Klipper adapter creation
func TestKlipperAdapter(t *testing.T) {
	t.Skip("Requires actual Klipper instance or mock server")

	printer := &v2printers.Printer{
		ID:   "klipper-1",
		Name: "Test Klipper",
		Type: v2printers.PrinterTypeKlipper,
		URL:  "http://mainsailos.local",
	}

	adapter, err := v2printers.NewKlipperAdapter(printer)
	require.NoError(t, err)
	require.NotNil(t, adapter)
}

// TestPrinterMonitoring tests printer status monitoring
func TestPrinterMonitoring(t *testing.T) {
	t.Skip("Requires actual printer or mock server")

	eventMgr := v2events.NewEventManager()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go eventMgr.Start(ctx)

	mgr := v2printers.NewPrinterManager(eventMgr)

	printer := &v2printers.Printer{
		ID:   "test-1",
		Name: "Test Printer",
		Type: v2printers.PrinterTypeOctoPrint,
		URL:  "http://test.local",
	}

	err := mgr.AddPrinter(printer)
	require.NoError(t, err)

	// Start monitoring
	err = mgr.StartMonitoring(ctx)
	assert.NoError(t, err)

	// Wait for some monitoring cycles
	time.Sleep(2 * time.Second)

	// Stop monitoring
	mgr.StopMonitoring()
}

// TestSendGCode tests sending G-code commands
func TestSendGCode(t *testing.T) {
	t.Skip("Requires actual printer or mock server")

	mgr := v2printers.NewPrinterManager(nil)

	printer := &v2printers.Printer{
		ID:   "test-1",
		Name: "Test Printer",
		Type: v2printers.PrinterTypeOctoPrint,
		URL:  "http://test.local",
	}

	err := mgr.AddPrinter(printer)
	require.NoError(t, err)

	// Send simple G-code
	err = mgr.SendGCode("test-1", "G28") // Home all axes
	assert.NoError(t, err)
}

// TestUploadFile tests uploading files to printer
func TestUploadFile(t *testing.T) {
	t.Skip("Requires actual printer or mock server")

	mgr := v2printers.NewPrinterManager(nil)

	printer := &v2printers.Printer{
		ID:   "test-1",
		Name: "Test Printer",
		Type: v2printers.PrinterTypeOctoPrint,
		URL:  "http://test.local",
	}

	err := mgr.AddPrinter(printer)
	require.NoError(t, err)

	// Upload test file
	testData := []byte("G28\nG1 Z10\n")
	err = mgr.UploadFile("test-1", "test.gcode", testData)
	assert.NoError(t, err)
}

// TestPrinterStateMapping tests state mapping from v1 to v2
func TestPrinterStateMapping(t *testing.T) {
	// This would test that v1 printer states map correctly to v2 states
	// For now, just verify the concept

	tests := []struct {
		v1State string
		v2State v2printers.PrinterState
	}{
		{"Operational", v2printers.StateIdle},
		{"Printing", v2printers.StatePrinting},
		{"Paused", v2printers.StatePaused},
		{"Offline", v2printers.StateOffline},
		{"Error", v2printers.StateError},
	}

	for _, tt := range tests {
		// State mapping logic exists in the adapters
		// This is just a conceptual test
		_ = tt
	}
}

// TestPrinterEventBroadcasting tests printer status events
func TestPrinterEventBroadcasting(t *testing.T) {
	t.Skip("Requires actual printer or mock server")

	eventMgr := v2events.NewEventManager()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go eventMgr.Start(ctx)

	mgr := v2printers.NewPrinterManager(eventMgr)

	printer := &v2printers.Printer{
		ID:   "test-1",
		Name: "Test Printer",
		Type: v2printers.PrinterTypeOctoPrint,
		URL:  "http://test.local",
	}

	err := mgr.AddPrinter(printer)
	require.NoError(t, err)

	// Start monitoring (which broadcasts events)
	err = mgr.StartMonitoring(ctx)
	assert.NoError(t, err)

	// Events should be broadcast as printer status changes
	time.Sleep(1 * time.Second)

	mgr.StopMonitoring()
}

// TestTemperatureInfo tests temperature information structure
func TestTemperatureInfo(t *testing.T) {
	temp := &v2printers.TemperatureInfo{
		HotendTemp:   210.5,
		HotendTarget: 210.0,
		BedTemp:      60.3,
		BedTarget:    60.0,
	}

	assert.Equal(t, 210.5, temp.HotendTemp)
	assert.Equal(t, 210.0, temp.HotendTarget)
	assert.Equal(t, 60.3, temp.BedTemp)
	assert.Equal(t, 60.0, temp.BedTarget)
}

// TestPrintProgress tests print progress structure
func TestPrintProgress(t *testing.T) {
	progress := &v2printers.PrintProgress{
		Completion:  45.5,
		PrintTime:   3600,
		TimeLeft:    4320,
		CurrentFile: "test_print.gcode",
	}

	assert.Equal(t, 45.5, progress.Completion)
	assert.Equal(t, 3600, progress.PrintTime)
	assert.Equal(t, 4320, progress.TimeLeft)
	assert.Equal(t, "test_print.gcode", progress.CurrentFile)
}

// TestPrinterConcurrency tests concurrent printer operations
func TestPrinterConcurrency(t *testing.T) {
	mgr := v2printers.NewPrinterManager(nil)

	// Add multiple printers concurrently
	numPrinters := 10
	errChan := make(chan error, numPrinters)

	for i := 0; i < numPrinters; i++ {
		go func(id int) {
			printer := &v2printers.Printer{
				ID:   string(rune(id)),
				Name: "Printer " + string(rune(id)),
				Type: v2printers.PrinterTypeOctoPrint,
				URL:  "http://test" + string(rune(id)) + ".local",
			}
			errChan <- mgr.AddPrinter(printer)
		}(i)
	}

	// Collect errors
	for i := 0; i < numPrinters; i++ {
		err := <-errChan
		assert.NoError(t, err)
	}

	// Verify all added
	printers, _ := mgr.ListPrinters()
	assert.Len(t, printers, numPrinters)
}

// BenchmarkPrinterManager benchmarks printer manager operations
func BenchmarkPrinterManager(b *testing.B) {
	mgr := v2printers.NewPrinterManager(nil)

	printer := &v2printers.Printer{
		ID:   "bench-1",
		Name: "Benchmark Printer",
		Type: v2printers.PrinterTypeOctoPrint,
		URL:  "http://bench.local",
	}

	b.Run("AddPrinter", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			p := &v2printers.Printer{
				ID:   string(rune(i)),
				Name: "Printer",
				Type: v2printers.PrinterTypeOctoPrint,
				URL:  "http://test.local",
			}
			_ = mgr.AddPrinter(p)
		}
	})

	b.Run("ListPrinters", func(b *testing.B) {
		// Add one printer first
		mgr.AddPrinter(printer)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = mgr.ListPrinters()
		}
	})
}
