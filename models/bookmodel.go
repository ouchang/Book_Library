package models

import (
	"fmt"
	"gobooklibrary/config"
)

type Book struct {
	Id              uint   `json:"id"`
	Title           string `json:"title" binding:"required"`
	Author          string `json:"author" binding:"required"`
	PublicationYear int    `json:"publication_year" binding:"required"`
	CategoryId      int    `json:"category_id" binding:"required"`
}

type NewBook struct {
	Title           string `json:"title" binding:"required"`
	Author          string `json:"author" binding:"required"`
	PublicationYear int    `json:"publication_year" binding:"required"`
	Category        string `json:"category" binding:"required"`
}

func (b *Book) TableName() string {
	return "Books"
}

// GetAllBooks Fetch all book data
func GetAllBooks(books *[]Book) (err error) {
	if err = config.ORMDB.Find(books).Error; err != nil {
		return err
	}
	return nil
}

func GetBooksByTitle(books *[]Book, title string) (err error) {
	if err = config.ORMDB.Where("title LIKE ?", "%"+title+"%").Find(&books).Error; err != nil {
		return err
	}
	return nil
}

func AddBook(book NewBook) (err error) {
	var addBookCode int

	bookAdd, err := config.DB.Query("call addBook(?, ?, ?, ?)", book.Title, book.Author, book.PublicationYear, book.Category)
	if err != nil {
		return err
	}

	for bookAdd.Next() {
		err2 := bookAdd.Scan(&addBookCode)
		if err2 != nil {
			return err2
		}
	}

	defer bookAdd.Close()

	if addBookCode == 0 {
		return fmt.Errorf("mysql procedure addBook failed")
	}

	return nil
}
