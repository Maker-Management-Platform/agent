package klipper

import (
	"github.com/eduardooliveira/stLib/core/events"
	printerEntities "github.com/eduardooliveira/stLib/core/printers/entities"
	"github.com/gorilla/websocket"
)

type KlipperPrinter struct {
	config                *printerEntities.Config
	state                 *printerEntities.State
	bedState              *printerEntities.TemperatureStatus
	hotEndState           []*printerEntities.TemperatureStatus
	jobState              *printerEntities.JobStatus
	ws                    *websocket.Conn
	bedChangeListeners    map[chan *printerEntities.TemperatureStatus]struct{}
	hotEndChangeListeners map[chan []*printerEntities.TemperatureStatus]struct{}
	jobChangeListeners    map[chan *printerEntities.JobStatus]struct{}
	deltaListeners        map[chan []*events.Message]struct{}
}

type MoonRakerResponse struct {
	Result *Result `json:"result"`
}
type Result struct {
	KlippyConnected           bool          `json:"klippy_connected"`
	KlippyState               string        `json:"klippy_state"`
	Components                []string      `json:"components"`
	FailedComponents          []interface{} `json:"failed_components"`
	RegisteredDirectories     []string      `json:"registered_directories"`
	Warnings                  []interface{} `json:"warnings"`
	WebsocketCount            int           `json:"websocket_count"`
	MoonrakerVersion          string        `json:"moonraker_version"`
	MissingKlippyRequirements []interface{} `json:"missing_klippy_requirements"`
	APIVersion                []int         `json:"api_version"`
	APIVersionString          string        `json:"api_version_string"`
}

type statusUpdate struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
}
type result struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Status    map[string]map[string]any `json:"status"`
		EventTime float64                   `json:"eventtime"`
	} `json:"result"`
	ID int `json:"id"`
}

type thermalStatus struct {
	Temperature float64 `json:"temperature"`
	Target      float64 `json:"target"`
}

type printStatsStatus struct {
	TotalDuration float64 `json:"total_duration"`
	FilamentUsed  float64 `json:"filament_used"`
}

type displayStatus struct {
	Message  string  `json:"message"`
	Progress float64 `json:"progress"`
}
