package database

import (
	"errors"

	"github.com/eduardooliveira/stLib/core/entities"
	"gorm.io/gorm"
)

func initPrintJob() error {
	return DB.AutoMigrate(&entities.PrintJob{})
}

func InsertPrintJob(p *entities.PrintJob) error {
	//return DB.Transaction(func(tx *gorm.DB) error {
	pos, err := lastPrintJobInQueue(DB)
	if err != nil {
		return err
	}
	p.Position = pos + 1
	if err := DB.Create(p).Error; err != nil {
		return err
	}
	return nil
	//})
}

func lastPrintJobInQueue(tx *gorm.DB) (int, error) {
	var p *entities.PrintJob
	if err := tx.Order("position desc").Take(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		} else {
			return 0, err
		}
	}
	return p.Position, nil
}

func GetPrintJobs(states []string) (rtn []*entities.PrintJob, err error) {
	if err = DB.Debug().Where("state in ?", states).Order("position asc").Preload("Slice").Find(&rtn).Error; err != nil {
		return nil, err
	}
	return rtn, nil
}
