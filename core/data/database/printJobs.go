package database

import (
	"errors"

	"github.com/eduardooliveira/stLib/core/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func initPrintJob() error {
	return DB.AutoMigrate(&entities.PrintJob{})
}

func InsertPrintJob(p *entities.PrintJob) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		pos, err := lastPrintJobInQueue(DB)
		if err != nil {
			return err
		}
		p.Position = pos + 1
		if err := DB.Create(p).Error; err != nil {
			return err
		}
		return nil
	})
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

func GetPrintJobs() (rtn []*entities.PrintJob, err error) {
	if err = DB.Debug().Order("position asc").Preload("Slice").Find(&rtn).Error; err != nil {
		return nil, err
	}
	return rtn, nil
}

func Move(id string, pos int) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		p, err := GetPrintJobByUUID(id)
		if err != nil {
			return err
		}

		if p.Position > pos {
			for i := p.Position; i >= pos; i-- {
				if err := tx.Debug().Model(&entities.PrintJob{}).Where("position = ?", i).Update("position", i+1).Error; err != nil {
					return err
				}
			}
		} else {
			for i := p.Position; i <= pos; i++ {
				if err := tx.Debug().Model(&entities.PrintJob{}).Where("position = ?", i).Update("position", i-1).Error; err != nil {
					return err
				}
			}
		}

		if err := tx.Debug().Model(p).Update("position", pos).Error; err != nil {
			return err
		}
		return nil
	})
}

func ChangePrintJobState(uuid string, state string) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var affected []*entities.PrintJob
		if err := tx.Debug().Clauses(clause.Returning{}).Model(&affected).Where("uuid = ?", uuid).Update("state", state).Error; err != nil {
			return err
		}

		if affected[0].Position > 0 {
			var initalPosition = affected[0].Position
			if err := tx.Debug().Model(affected[0]).Update("position", -1).Error; err != nil {
				return err
			}
			lastPosition, err := lastPrintJobInQueue(tx)
			if err != nil {
				return err
			}
			for i := initalPosition; i <= lastPosition; i++ {
				if err := tx.Debug().Model(&entities.PrintJob{}).Where("position = ?", i).Update("position", i-1).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func SwapPrintJobs(ps1, ps2 string) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		p1, err := GetPrintJobByUUID(ps1)
		if err != nil {
			return err
		}
		p2, err := GetPrintJobByUUID(ps2)
		if err != nil {
			return err
		}

		if err := tx.Model(p1).Update("position", p2.Position).Error; err != nil {
			return err
		}
		if err := tx.Model(p2).Update("position", p1.Position).Error; err != nil {
			return err
		}
		return nil
	})
}

func GetPrintJobByUUID(uuid string) (*entities.PrintJob, error) {
	var p *entities.PrintJob
	if err := DB.Debug().Where("uuid = ?", uuid).Take(&p).Error; err != nil {
		return nil, err
	}
	return p, nil
}
