package entities

import "github.com/google/uuid"

const (
	PrintJobStatusQueued   = "queued"
	PrintJobStatusPrinting = "printing"
	PrintJobStatusDone     = "done"
)

type PrintJob struct {
	UUID string `json:"uuid"`
	//Tags   []*Tag        `json:"tags" gorm:"many2many:project_tags"`
	Slice    *ProjectAsset `json:"slice" gorm:"foreignKey:SliceId"`
	SliceId  int           `json:"-"`
	Position int           `json:"position"`
	Status   string        `json:"status"`
	Result   string        `json:"result"`
}

func NewPrintJob(sliceAsset *ProjectAsset) *PrintJob {
	return &PrintJob{
		UUID:   uuid.New().String(),
		Status: PrintJobStatusQueued,
	}
}
