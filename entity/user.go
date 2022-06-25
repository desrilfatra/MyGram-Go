package entity

import "time"

type User struct {
	Id        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"createdat"`
	UpdatedAt time.Time `json:"updatedat"`
}

type ResponseRegister struct {
	Age      int    `json:"age"`
	Email    string `json:"email"`
	Id       int    `json:"id"`
	Username string `json:"username"`
}
