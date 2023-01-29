package controllers

import (
	"gobooklibrary/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetStatuses(c *gin.Context) {
	var statuses []models.Status
	err := models.GetAllStatus(&statuses)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, statuses)
	}
}

func GetAvaliable(c *gin.Context) {
	var statuses []models.Status
	err := models.GetAvaliable(&statuses)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, statuses)
	}
}
