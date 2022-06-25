package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"mygram-go/entity"
	"mygram-go/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	db *sql.DB
}

// UserHandler implements UserHandlerIF
func (h *UserHandler) UserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	switch r.Method {
	case http.MethodGet:
		//users/{id}
		if id != "" { // get by id
			h.getUsersByIDHandler(w, r, id)
		} else { // get all
			//users
			h.getUsersHandler(w, r)
		}
	case http.MethodPost:
		//users
		h.createUsersHandler(w, r)
	case http.MethodPut:
		//users/{id}
		h.updateUserHandler(w, r, id)
	case http.MethodDelete:
		//users/{id}
		h.deleteUserHandler(w, r, id)
	}
}

type RegisterHandler struct {
	db *sql.DB
}

// RegisterUser implements RegisterHandlerIF
func (h *RegisterHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var newUser entity.User
	json.NewDecoder(r.Body).Decode(&newUser)
	newPassword := []byte(newUser.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(newPassword, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	newUser.Password = string(hashedPassword)

	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()
	sqlStatment := `INSERT INTO public.users
	(username,email,password,age,createdat,updatedat)
	values ($1,$2,$3,$4,$5,$6) Returning id` //sesuai dengan nama table
	fmt.Println("tess")
	err = h.db.QueryRow(sqlStatment,
		newUser.Username,
		newUser.Email,
		newUser.Password,
		newUser.Age,
		newUser.CreatedAt,
		newUser.UpdatedAt,
	).Scan(&newUser.Id)
	fmt.Println(newUser.Id)
	response_Register := entity.ResponseRegister{
		Age:      newUser.Age,
		Email:    newUser.Email,
		Id:       newUser.Id,
		Username: newUser.Username,
	}
	jsonData, _ := json.Marshal(&response_Register)
	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonData)
}

type LoginHandler struct {
	db *sql.DB
}

// LoginUser implements LoginHandlerIF
func (h *LoginHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var newUser entity.User
		var validasiUser *entity.User

		json.NewDecoder(r.Body).Decode(&newUser)
		fmt.Println(r.Body)
		tempPassword := newUser.Password
		newPassword := []byte(newUser.Password)
		_, err := bcrypt.GenerateFromPassword(newPassword, bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		validasiUser = &newUser
		serv := service.NewUserService()

		// newUser.Password = string(hashedPassword)
		// fmt.Println(newUser.Password)
		sqlStatment := `select * from public.users where email = $1`

		err = h.db.QueryRow(sqlStatment, newUser.Email).
			Scan(&newUser.Id, &newUser.Username, &newUser.Email, &newUser.Password, &newUser.Age, &newUser.CreatedAt, &newUser.UpdatedAt)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(newUser)
		validasiUser, err = serv.Login(validasiUser, tempPassword)
		if err != nil {
			fmt.Println(err)
			w.Write([]byte(fmt.Sprint(err)))
		} else {
			var token entity.Tokens
			token.TokenJwt = serv.GetToken(uint(newUser.Id), newUser.Email, newUser.Password)
			jsonData, _ := json.Marshal(&token)
			w.Header().Add("Content-Type", "application/json")
			w.Write(jsonData)
		}

	} else {
		fmt.Println("ERORRRR")
	}
}

type UserHandlerIF interface {
	UserHandler(w http.ResponseWriter, r *http.Request)
}

type RegisterHandlerIF interface {
	RegisterUser(w http.ResponseWriter, r *http.Request)
}

type LoginHandlerIF interface {
	LoginUser(w http.ResponseWriter, r *http.Request)
}

func NewUserHandler(db *sql.DB) UserHandlerIF {
	return &UserHandler{db: db}
}

func NewRegisterHandler(db *sql.DB) RegisterHandlerIF {
	return &RegisterHandler{db: db}
}

func UserLoginHandler(db *sql.DB) LoginHandlerIF {
	return &LoginHandler{db: db}
}

func (h *UserHandler) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	users := []*entity.User{}
	sqlStatment := `SELECT  * from users` //sesuai dengan nama table

	rows, err := h.db.Query(sqlStatment)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var user entity.User
		if serr := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Age, &user.CreatedAt, &user.UpdatedAt); serr != nil {
			fmt.Println("Scan error", serr)
		}
		users = append(users, &user)
	}
	jsonData, _ := json.Marshal(&users)
	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonData)
}

func (h *UserHandler) getUsersByIDHandler(w http.ResponseWriter, r *http.Request, id string) {
	users := []*entity.User{}
	sqlStatment := `SELECT  * from users where id = $1` //sesuai dengan nama table
	rows, err := h.db.Query(sqlStatment)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var user entity.User
		if serr := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Age, &user.CreatedAt, &user.UpdatedAt); serr != nil {
			fmt.Println("Scan error", serr)
		}
		users = append(users, &user)
	}
	jsonData, _ := json.Marshal(&users)
	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonData)
}

func (h *UserHandler) createUsersHandler(w http.ResponseWriter, r *http.Request) {

	var newUser entity.User
	json.NewDecoder(r.Body).Decode(&newUser)
	sqlStatment := `insert into users
	(username,email,password,age,createdat,updatedat)
	values ($1,$2,$3,$4,$5,$5)` //sesuai dengan nama table
	res, err := h.db.Exec(sqlStatment,
		newUser.Username,
		newUser.Email,
		newUser.Password,
		newUser.Age,
		time.Now(),
	)

	if err != nil {
		panic(err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}

	w.Write([]byte(fmt.Sprint("User  update ", count)))
	return
}

func (h *UserHandler) updateUserHandler(w http.ResponseWriter, r *http.Request, id string) {
	if id != "" { // get by id
		var newUser entity.User
		json.NewDecoder(r.Body).Decode(&newUser)
		newUser.CreatedAt = time.Now()
		newUser.UpdatedAt = time.Now()
		sqlstatment := `
		update users set username = $1, email = $2, password = $3, createdat = $4, updatedat = $5 
		where id = $6`

		res, err := h.db.Exec(sqlstatment,
			newUser.Username,
			newUser.Email,
			newUser.Password,
			newUser.CreatedAt,
			newUser.UpdatedAt,
		)

		if err != nil {
			panic(err)
		}
		count, err := res.RowsAffected()
		if err != nil {
			panic(err)
		}

		w.Write([]byte(fmt.Sprint("User  update ", count)))
		return
	}
}

func (h *UserHandler) deleteUserHandler(w http.ResponseWriter, r *http.Request, id string) {
	sqlstament := `DELETE from users where id = $1;`
	if idInt, err := strconv.Atoi(id); err == nil {
		res, err := h.db.Exec(sqlstament, idInt)
		if err != nil {
			panic(err)
		}
		count, err := res.RowsAffected()
		if err != nil {
			panic(err)
		}
		w.Write([]byte(fmt.Sprint("Delete user rows ", count)))
		return
	}

}
