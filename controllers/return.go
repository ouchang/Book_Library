package controllers

import (
	"gobooklibrary/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReturnBook(c *gin.Context) {
	var returned models.Return
	var err error

	if err = c.BindJSON(&returned); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if err = models.ReturnBook(returned); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.IndentedJSON(http.StatusCreated, returned)
		}
	}
}
