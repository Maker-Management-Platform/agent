package octorpint

import models "github.com/eduardooliveira/stLib/core/entities"

type OctoPrintPrinter struct {
	*models.Printer
}

type OctoPrintResponse struct {
	APIVersion string `json:"version"`
	Safemode   string `json:"safemode"`
}
