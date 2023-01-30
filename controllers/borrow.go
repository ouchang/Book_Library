package controllers

import (
	"gobooklibrary/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBorrows(c *gin.Context) {
	var borrows []models.Borrow
	var err error

	user_login, user_login_ok := c.GetQuery("login")

	if user_login_ok {
		err = models.GetUsersBorrows(&borrows, user_login)
	} else {
		err = models.GetAllBorrows(&borrows)
	}

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, borrows)
	}
}

func GetMyBorrows(c *gin.Context) {
	var borrows []models.Borrow
	var err error

	user_login, user_login_ok := c.GetQuery("login")

	if user_login_ok {
		err = models.GetUsersBorrows(&borrows, user_login)
	}

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

func ReturnBook(c *gin.Context) {
	var returned models.ReturnRenew
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

func RenewBook(c *gin.Context) {
	var renewed models.ReturnRenew
	var err error

	if err = c.BindJSON(&renewed); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if err = models.RenewBook(renewed); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.IndentedJSON(http.StatusCreated, renewed)
		}
	}
}
