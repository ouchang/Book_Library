package models

import "gobooklibrary/config"

type Status struct {
	BookId uint   `json:"book_id"`
	Status string ` json:"status"`
}

func (s *Status) TableName() string {
	return "BookStatus"
}

func GetAllStatus(statuses *[]Status) (err error) {
	if err = config.ORMDB.Find(statuses).Error; err != nil {
		return err
	}

	return nil
}

func GetAvaliable(statuses *[]Status) (err error) {
	if err = config.ORMDB.Where("status = ?", "avaliable").Find(statuses).Error; err != nil {
		return err
	}

	return nil
}
