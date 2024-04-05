package entities

type PrintQueue struct {
	Jobs []*PrintJob `json:"jobs"`
}
