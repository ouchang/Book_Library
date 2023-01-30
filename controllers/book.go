package controllers

import (
	"gobooklibrary/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBooks(c *gin.Context) {
	var books []models.Book
	var err error

	title, title_ok := c.GetQuery("title")
	author, author_ok := c.GetQuery("author")

	if title_ok {
		err = models.GetBooksByTitle(&books, title)
	} else if author_ok {
		err = models.GetBooksByAuthor(&books, author)
	} else {
		err = models.GetAllBooks(&books)
	}

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, books)
	}
}

func AddBook(c *gin.Context) {
	var book models.NewBook
	var err error

	if err = c.BindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if err = models.AddBook(book); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.IndentedJSON(http.StatusCreated, book)
		}
	}
}
