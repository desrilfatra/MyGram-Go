package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
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
	sqlQuery := `INSERT INTO users
				(username,email,password,age,created_at,updated_at) 
				values (?,?,?,?,?,?)`
	fmt.Println("tess")
	res, err := db.Exec(sqlQuery,
		newUser.Username,
		newUser.Email,
		newUser.Password,
		newUser.Age,
		newUser.CreatedAt,
		newUser.UpdatedAt,
	)
	if err != nil {
		panic(err.Error())
	}

	lastId, err := res.LastInsertId()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("The last inserted row id: %d\n", lastId)
	if err != nil {
		panic(err)
	} else {
		response_Register := entity.ResponseRegister{
			Age:      newUser.Age,
			Email:    newUser.Email,
			Id:       int(lastId),
			Username: newUser.Username,
		}
		return response_Register
	}
}

func UserLoginRepository(db *sql.DB, user entity.User) (entity.User, error) {
	sqlQuery := `select u.id, u.username, u.email, u.password, u.age,
				u.created_at, u.updated_at from users as u where email = ?`
	res, err := db.Query(sqlQuery, user.Email)
	defer res.Close()

	if err != nil {
		log.Fatal(err)
	}
	for res.Next() {
		err := res.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Age, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err != nil {
		return entity.User{}, errors.New("username cannot be empty")
	}
	return user, nil
}

func UserPutRepository(db *sql.DB, NewUser entity.User, id string) entity.ResponseUpdate {

	sqlQuery := `
		UPDATE users set username = ?, email= ?, updated_at = ? 
		where id = ?`

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
		u.created_at, u.updated_at from users as u  where id= ?`

	res, err := db.Query(sqlQuery1, id)
	defer res.Close()

	if err != nil {
		log.Fatal(err)
	}
	for res.Next() {
		err := res.Scan(&NewUser.Id, &NewUser.Username, &NewUser.Email,
			&NewUser.Password, &NewUser.Age, &NewUser.CreatedAt, &NewUser.UpdatedAt)
		if err != nil {
			log.Fatal(err)
		}
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
	sqlQuery := `DELETE FROM users where id = ?`
	_, err := db.Exec(sqlQuery, newUser.Id)
	if err != nil {
		panic(err)
	}
	responseDel := entity.Message{
		Message: "Your account has been successfully deleted",
	}
	return responseDel
}
