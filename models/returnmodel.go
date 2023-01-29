package models

import (
	"fmt"
	"gobooklibrary/config"
)

type Return struct {
	UserId uint `json:"user_id" binding:"required"`
	BookId uint `json:"book_id" binding:"required"`
}

func ReturnBook(returned Return) (err error) {
	var registerCode int
	insert, err := config.DB.Query("call returnBook(?, ?)", returned.UserId, returned.BookId)

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
		return fmt.Errorf("mysql procedure returnBook failed")
	}

	return nil
}
