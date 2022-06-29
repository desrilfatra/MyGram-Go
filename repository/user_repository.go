package repository

import (
	"database/sql"
	"fmt"
	"mygram-go/entity"
	"time"
)

var (
	db  *sql.DB
	err error
)

func UserRegisterRepository(db *sql.DB, newUser entity.User) entity.ResponseRegister {
	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()
	sqlQuery := `INSERT INTO public.users
				(username,email,password,age,created_at,updated_at)
				values ($1,$2,$3,$4,$5,$6) Returning id`
	fmt.Println("tess")
	err = db.QueryRow(sqlQuery,
		newUser.Username,
		newUser.Email,
		newUser.Password,
		newUser.Age,
		newUser.CreatedAt,
		newUser.UpdatedAt,
	).Scan(&newUser.Id)
	if err != nil {
		panic(err)
	} else {
		response_Register := entity.ResponseRegister{
			Age:      newUser.Age,
			Email:    newUser.Email,
			Id:       newUser.Id,
			Username: newUser.Username,
		}
		return response_Register
	}

}
