package models

import (
	"fmt"
	"gobooklibrary/config"
)

type Copy struct {
	Id     int    `json:"id"`
	BookId uint   `json:"book_id" binding:"required"`
	Status string ` json:"status"`
}

type ISBN struct {
	ISBN string `json:"isbn" binding:"required"`
}

type CopyId struct {
	Id int `json:"id"`
}

func (c *Copy) TableName() string {
	return "BookCopies"
}

func GetAllCopy(copies *[]Copy) (err error) {
	if err = config.ORMDB.Find(copies).Error; err != nil {
		return err
	}

	return nil
}

func GetAvaliable(copies *[]Copy) (err error) {
	if err = config.ORMDB.Where("status = ?", "avaliable").Find(copies).Error; err != nil {
		return err
	}

	return nil
}

func AddBookCopy(isbn ISBN) (err error) {
	var registerCode int
	insert, err := config.DB.Query("call addBookCopy(?)", isbn.ISBN)

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
		return fmt.Errorf("mysql procedure addBookCopy failed")
	}

	return nil
}

func DeleteBookCopy(copy_id CopyId) (err error) {
	var registerCode int
	insert, err := config.DB.Query("call deleteBookCopy(?)", copy_id.Id)

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
		return fmt.Errorf("mysql procedure deleteBookCopy failed")
	}

	return nil
}
