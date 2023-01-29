package models

type ISBN struct {
	ISBN string `json:"isbn" binding:"required"`
}
