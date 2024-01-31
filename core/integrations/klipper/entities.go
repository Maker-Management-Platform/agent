package klipper

import "github.com/eduardooliveira/stLib/core/models"

type KlipperPrinter struct {
	*models.Printer
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
		Status    map[string]any `json:"status"`
		EventTime float64        `json:"eventtime"`
	} `json:"result"`
	ID int `json:"id"`
}
