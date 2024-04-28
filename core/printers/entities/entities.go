package entities

import (
	"github.com/eduardooliveira/stLib/core/entities"
	"github.com/eduardooliveira/stLib/core/events"
)

type JobStatus struct {
	Name          string  `json:"-"`
	Progress      float64 `json:"progress,omitempty"`
	Status        string  `json:"status,omitempty"`
	Message       string  `json:"message,omitempty"`
	FileName      string  `json:"fileName,omitempty"`
	TotalDuration float64 `json:"totalDuration,omitempty"`
}

type TemperatureStatus struct {
	Name        string  `json:"-"`
	Temperature float64 `json:"temperature,omitempty"`
	Target      float64 `json:"target,omitempty"`
	Power       float64 `json:"power,omitempty"`
}

type Config struct {
	UUID      string `json:"uuid" toml:"uuid" form:"uuid" query:"uuid"`
	Name      string `json:"name" toml:"name" form:"name" query:"name"`
	Type      string `json:"type" toml:"type" form:"type" query:"type"`
	Address   string `json:"address" toml:"address" form:"address" query:"address"`
	CameraUrl string `json:"camera_url" toml:"camera_url" form:"camera_url" query:"camera_url"`
	ApiKey    string `json:"apiKey" toml:"apiKey" form:"apiKey" query:"apiKey"`
}

type State struct {
	ConnectionStatus string `json:"connectionStatus"`
	Version          string `json:"version"`
}

type Printer interface {
	GetConfig() *Config
	GetState() *State
	UploadFile(asset *entities.ProjectAsset) error
	TestConnection() error
	OnHotEndChange(out chan []*TemperatureStatus)
	OnBedChange(out chan *TemperatureStatus)
	OnJobChange(out chan *JobStatus)
	OnDeltaChange(out chan []*events.Message)
}
