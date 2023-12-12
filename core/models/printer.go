package models

type Printer struct {
	UUID    string `json:"uuid" toml:"uuid" form:"uuid" query:"uuid"`
	Name    string `json:"name" toml:"name" form:"name" query:"name"`
	Type    string `json:"type" toml:"type" form:"type" query:"type"`
	Address string `json:"address" toml:"address" form:"address" query:"address"`
}
