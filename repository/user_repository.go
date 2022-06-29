package repository

import (
	"database/sql"
	"errors"
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

func UserLoginRepository(db *sql.DB, user entity.User) (entity.User, error) {
	sqlQuery := `select u.id, u.username, u.email, u.password, u.age,
				u.created_at, u.updated_at from public.users as u  where email= $1`
	err = db.QueryRow(sqlQuery, user.Email).
		Scan(&user.Id, &user.Username, &user.Email, &user.Password,
			&user.Age, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return entity.User{}, errors.New("username cannot be empty")
	}
	return user, nil
}
