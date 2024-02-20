package klipper

import "github.com/eduardooliveira/stLib/core/entities"

func ConnectionStatus(printer *entities.Printer) error {
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

func UploadFile(printer *entities.Printer, asset *entities.ProjectAsset) error {
	kp := &KlipperPrinter{printer}

	return kp.ServerFilesUpload(asset)
}
