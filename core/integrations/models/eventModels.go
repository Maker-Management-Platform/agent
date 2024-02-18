package models

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
