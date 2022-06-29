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

func UserPutRepository(db *sql.DB, NewUser entity.User, id string) entity.ResponseUpdate {

	sqlQuery := `
		UPDATE public.users set username = $1, email= $2, updated_at = $3 
		where id = $4`

	_, err := db.Exec(sqlQuery,
		NewUser.Username,
		NewUser.Email,
		time.Now(),
		id,
	)
	if err != nil {
		fmt.Println("error update")
		panic(err)

	}
	sqlQuery1 := `select u.id, u.username, u.email, u.password, u.age,
		u.created_at, u.updated_at from public.users as u  where id= $1`
	err = db.QueryRow(sqlQuery1, id).
		Scan(&NewUser.Id, &NewUser.Username, &NewUser.Email, &NewUser.Password,
			&NewUser.Age, &NewUser.CreatedAt, &NewUser.UpdatedAt)

	if err != nil {
		panic(err)
	}

	fmt.Println(NewUser)

	responseUpdateUser := entity.ResponseUpdate{
		Id:        NewUser.Id,
		Email:     NewUser.Email,
		Username:  NewUser.Username,
		Age:       NewUser.Age,
		UpdatedAt: time.Now(),
	}
	return responseUpdateUser
}

func UserDeleteRepository(db *sql.DB, newUser *entity.User) entity.Message {
	sqlQuery := `DELETE FROM public.users where id = $1`
	_, err := db.Exec(sqlQuery, newUser.Id)
	if err != nil {
		panic(err)
	}
	responseDel := entity.Message{
		Message: "Your account has been successfully deleted",
	}
	return responseDel
}
