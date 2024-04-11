package entities

import "github.com/google/uuid"

const (
	PrintJobStateQueued   = "queued"
	PrintJobStatePrinting = "printing"
	PrintJobStateDone     = "done"
)

type PrintJob struct {
	UUID string `json:"uuid" gorm:"primaryKey"`
	//Tags   []*Tag        `json:"tags" gorm:"many2many:project_tags"`
	Slice    *ProjectAsset `json:"slice" gorm:"references:ID"`
	SliceId  string        `json:"sliceId"`
	Position int           `json:"position"`
	State    string        `json:"state"`
	Result   string        `json:"result"`
}

func NewPrintJob(sliceAsset *ProjectAsset) *PrintJob {
	return &PrintJob{
		UUID:    uuid.New().String(),
		State:   PrintJobStateQueued,
		Slice:   sliceAsset,
		SliceId: sliceAsset.ID,
	}
}
