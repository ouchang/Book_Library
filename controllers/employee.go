package controllers

import (
	"gobooklibrary/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetEmployees(c *gin.Context) {
	var employees []models.Employee
	err := models.GetAllEmployees(&employees)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, employees)
	}
}

func AddEmployee(c *gin.Context) {
	var employee models.Employee
	var err error

	if err = c.BindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if err = models.AddEmployee(employee); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.IndentedJSON(http.StatusCreated, employee)
		}
	}
}

func DeleteEmployee(c *gin.Context) {
	var employee_id models.EmployeeId
	var err error

	if err = c.BindJSON(&employee_id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if err = models.DeleteEmployee(employee_id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.IndentedJSON(http.StatusCreated, employee_id)
		}
	}
}
