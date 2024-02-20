package entities

import "github.com/google/uuid"

type Printer struct {
	UUID      string `json:"uuid" toml:"uuid" form:"uuid" query:"uuid"`
	Name      string `json:"name" toml:"name" form:"name" query:"name"`
	Type      string `json:"type" toml:"type" form:"type" query:"type"`
	Address   string `json:"address" toml:"address" form:"address" query:"address"`
	CameraUrl string `json:"camera_url" toml:"camera_url" form:"camera_url" query:"camera_url"`
	Status    string `json:"status" toml:"status" form:"status" query:"status"`
	State     string `json:"state" toml:"state" form:"state" query:"state"`
	Version   string `json:"version" toml:"version" form:"version" query:"version"`
}

type PrinterStatus struct {
	Name  string `json:"name"`
	State any    `json:"state"`
}

func NewPrinter() *Printer {
	printer := &Printer{
		UUID: uuid.New().String(),
	}
	return printer
}
