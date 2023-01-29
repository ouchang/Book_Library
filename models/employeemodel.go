package models

import (
	"fmt"
	"gobooklibrary/config"
)

const (
	EmployeeAdmin     string = "admin"
	EmployeeLibrarian string = "librarian"
)

type Employee struct {
	Id           uint   `json:"id"`
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
	Login        string `json:"login" binding:"required"`
	Password     string `json:"password" binding:"required"`
	EmployeeType string `json:"employee_type" binding:"required"`
}

func (e *Employee) TableName() string {
	return "Employees"
}

func GetEmployees(employees *[]Employee) (err error) {
	if err = config.ORMDB.Find(employees).Error; err != nil {
		return err
	}

	return nil
}

func LoginEmployee(loginemployee LoginAuth, kind string) (err error) {
	var loginEmployeeCode int

	login, err := config.DB.Query("call logInLibrarian(?, ?, ?)", loginemployee.Login, loginemployee.Password, kind)

	if err != nil {
		return err
	}

	for login.Next() {
		err2 := login.Scan(&loginEmployeeCode)
		if err2 != nil {
			return err2
		}
	}

	defer login.Close()

	if loginEmployeeCode == 0 {
		return fmt.Errorf("mysql procedure logInLibrarian failed")
	}

	return nil
}
