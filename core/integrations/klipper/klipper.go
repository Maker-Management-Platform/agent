package klipper

import "github.com/eduardooliveira/stLib/core/models"

func ConnectionStatus(printer *models.Printer) error {
	kp := &KlipperPrinter{printer}

	r, err := kp.serverInfo()

	if err != nil {
		kp.Status = "disconnected"
		return err
	}

	kp.Version = r.APIVersionString
	kp.State = r.KlippyState
	kp.Status = "connected"

	return nil
}

func UploadFile(printer *models.Printer, asset *models.ProjectAsset) error {
	kp := &KlipperPrinter{printer}

	return kp.ServerFilesUpload(asset)
}
