package controllers

import (
	"gobooklibrary/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBooks(c *gin.Context) {
	var books []models.Book
	var err error

	title, ok := c.GetQuery("title")

	if ok {
		err = models.GetBooksByTitle(&books, title)
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
