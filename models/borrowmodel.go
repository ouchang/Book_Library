package models

import (
	"fmt"
	"gobooklibrary/config"
	"strconv"
)

type Borrow struct {
	Id           uint   `json:"id"`
	UserId       uint   `json:"user_id" binding:"required"`
	BookId       uint   `json:"book_id" binding:"required"`
	EmployeeId   uint   `json:"employee_id" binding:"required"`
	ReleaseDate  string `json:"release_date"`
	DueDate      string `json:"due_date"`
	ReturnedDate string `json:"returned_date"`
	Returned     bool   `json:"returned"`
}

type ReturnRenew struct {
	UserId uint `json:"user_id" binding:"required"`
	BookId uint `json:"book_id" binding:"required"`
}

func (b *Borrow) TableName() string {
	return "BorrowLog"
}

func GetAllBorrows(borrows *[]Borrow) (err error) {
	if err = config.ORMDB.Find(borrows).Error; err != nil {
		return err
	}

	return nil
}

func GetUsersBorrows(borrows *[]Borrow, login string) (err error) {
	var userId int
	select_query, err := config.DB.Query("SELECT id FROM Users WHERE login = ?", login)

	if err != nil {
		return err
	}

	for select_query.Next() {

		err2 := select_query.Scan(&userId)

		if err2 != nil {
			return err2
		}
	}

	if err = config.ORMDB.Where("user_id LIKE ?", "%"+strconv.Itoa(userId)+"%").Find(&borrows).Error; err != nil {
		return err
	}
	return nil
}

func AddBorrow(borrow Borrow) (err error) {
	var registerCode int
	insert, err := config.DB.Query("call borrowBook(?, ?, ?)", borrow.UserId, borrow.BookId, borrow.EmployeeId)

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
		return fmt.Errorf("mysql procedure borrowBook failed")
	}

	return nil
}

func ReturnBook(returned ReturnRenew) (err error) {
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

func RenewBook(renewed ReturnRenew) (err error) {
	var registerCode int
	insert, err := config.DB.Query("call renewBook(?, ?)", renewed.UserId, renewed.BookId)

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
		return fmt.Errorf("mysql procedure renewBook failed")
	}

	return nil
}
