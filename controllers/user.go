package controllers

import (
	"gobooklibrary/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	var users []models.User
	err := models.GetAllUsers(&users)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, users)
	}
}

func AddUser(c *gin.Context) {
	var user models.User
	var err error
	if err = c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if err = models.AddUser(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.IndentedJSON(http.StatusCreated, user)
		}
	}
}

func GetUserProfileInfo(c *gin.Context) {
	var profile models.ProfileInfo
	var err error
	if err = c.BindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if err = models.GetProfileInfo(profile); err != nil {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.JSON(http.StatusOK, profile)
		}
	}
}
