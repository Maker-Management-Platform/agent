package database

import (
	"github.com/eduardooliveira/stLib/core/entities"
	"gorm.io/gorm"
)

func initPrintJob() error {
	return DB.AutoMigrate(&entities.PrintJob{})
}

func InsertPrintJob(p *entities.PrintJob) error {

	tx := DB.Begin()

	pos, err := lastPrintJobInQueue(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	p.Position = pos + 1

	if err := DB.Create(p).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func lastPrintJobInQueue(tx *gorm.DB) (int, error) {
	var p entities.PrintJob
	if err := tx.Order("position asc").Take(p).Error; err != nil {
		return 0, err
	}
	return p.Position, nil
}
