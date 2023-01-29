package controllers

import (
	"gobooklibrary/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBorrows(c *gin.Context) {
	var borrows []models.Borrow
	err := models.GetAllBorrows(&borrows)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, borrows)
	}
}

func AddBorrow(c *gin.Context) {
	var borrow models.Borrow
	var err error

	if err = c.BindJSON(&borrow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if err = models.AddBorrow(borrow); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.IndentedJSON(http.StatusCreated, borrow)
		}
	}
}
