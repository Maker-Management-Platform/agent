package octorpint

import "github.com/eduardooliveira/stLib/core/models"

func ConnectionStatus(printer *models.Printer) error {
	op := &OctoPrintPrinter{printer}

	r, err := op.serverInfo()

	if err != nil {
		op.Status = "disconnected"
		return err
	}

	op.Version = r.APIVersion
	op.Status = "connected"

	return nil
}

func UploadFile(printer *models.Printer, asset *models.ProjectAsset) error {
	op := &OctoPrintPrinter{printer}

	return op.ServerFilesUpload(asset)
}
