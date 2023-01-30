package controllers

import (
	"gobooklibrary/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCopies(c *gin.Context) {
	var copies []models.Copy
	err := models.GetAllCopy(&copies)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, copies)
	}
}

func GetAvaliable(c *gin.Context) {
	var copies []models.Copy
	err := models.GetAvaliable(&copies)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, copies)
	}
}

func AddBookCopy(c *gin.Context) {
	var isbn models.ISBN
	var err error
	if err = c.BindJSON(&isbn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if err = models.AddBookCopy(isbn); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.IndentedJSON(http.StatusCreated, isbn)
		}
	}
}

func DeleteBookCopy(c *gin.Context) {
	var copy_id models.CopyId
	var err error
	if err = c.BindJSON(&copy_id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if err = models.DeleteBookCopy(copy_id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.IndentedJSON(http.StatusCreated, copy_id)
		}
	}
}
