package models

import (
	"fmt"
	"gobooklibrary/config"
)

type Category struct {
	Id       uint   `json:"idx"`
	Category string `json:"category" binding:"required"`
}

func (c *Category) TableName() string {
	return "BookCategories"
}

func GetAllCategories(categories *[]Category) (err error) {
	if err = config.ORMDB.Find(categories).Error; err != nil {
		return err
	}

	return nil
}

func AddCategory(category Category) (err error) {
	var registerCode int
	insert, err := config.DB.Query("call addCategory(?)", category.Category)

	if err != nil {
		return err
	}

	for insert.Next() {

		err2 := insert.Scan(&registerCode)

		if err2 != nil {
			return err2
		}
	}

	defer insert.Close()

	if registerCode == 0 {
		return fmt.Errorf("mysql procedure addCategory failed")
	}

	return nil
}
