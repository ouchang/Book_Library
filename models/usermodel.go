package models

import (
	"fmt"
	"gobooklibrary/config"
)

type User struct {
	Id          uint   `json:"id"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	CreatedYear int    `json:"created_year"`
	PhoneNumber int    `json:"phone_number" binding:"required"`
	Login       string `json:"login" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

func (u *User) TableName() string {
	return "Users"
}

func GetAllUsers(users *[]User) (err error) {
	if err = config.ORMDB.Find(users).Error; err != nil {
		return err
	}

	return nil
}

func LoginUser(loginuser LoginAuth) (err error) {
	var loginUserCode int
	login, err := config.DB.Query("call logInUser(?, ?)", loginuser.Login, loginuser.Password)

	if err != nil {
		return err
	}

	for login.Next() {
		err2 := login.Scan(&loginUserCode)
		if err2 != nil {
			return err2
		}
	}

	defer login.Close()

	if loginUserCode == 0 {
		return fmt.Errorf("mysql procedure logInUser failed")
	}

	return nil
}

func AddUser(user User) (err error) {
	var registerCode int
	insert, err := config.DB.Query("call registerUser(?, ?, ?, ?, ?)", user.FirstName, user.LastName, user.PhoneNumber, user.Login, user.Password)

	if err != nil {
		return err
	}

	for insert.Next() {

		err2 := insert.Scan(&registerCode)

		if err2 != nil {
			return err2
		}
	}

	defer insert.Close()

	if registerCode == 0 {
		return fmt.Errorf("mysql procedure registerUser failed")
	}

	return nil
}
