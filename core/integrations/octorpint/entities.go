package octorpint

import "github.com/eduardooliveira/stLib/core/models"

type OctoPrintPrinter struct {
	*models.Printer
}

type OctoPrintResponse struct {
	APIVersion string `json:"version"`
	Safemode   string `json:"safemode"`
}
