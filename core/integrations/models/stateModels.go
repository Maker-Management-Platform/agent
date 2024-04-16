package models

import "github.com/eduardooliveira/stLib/core/entities"

type PrinterState struct {
	Printer  *entities.Printer
	State    string
	Progress float64
	Duration int64 // use klipper print stats totalDuration
	JobUUID  string
}
