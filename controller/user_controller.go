package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"mygram-go/entity"
	"mygram-go/middleware"
	"mygram-go/repository"
	"mygram-go/service"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type UsersHandler struct {
	db *sql.DB
}

// UsersHandler implements UsersHandlerIF
func (h *UsersHandler) UsersHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	switch r.Method {
	case http.MethodPut:
		h.Updateusr(w, r, id)
	case http.MethodDelete:
		h.Delete(w, r)
	}
}

type RegisterHandler struct {
	db *sql.DB
}

// Register implements RegisterHandlerIF
func (h *RegisterHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var newUser entity.User
		var validate *entity.User
		json.NewDecoder(r.Body).Decode(&newUser)
		newPassword := []byte(newUser.Password)
		hashedPassword, err := bcrypt.GenerateFromPassword(newPassword, bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		validate = &newUser
		serv := service.NewUserService()
		validate, err = serv.Register(validate)
		if err != nil {
			w.Write([]byte(fmt.Sprint(err)))

		} else {
			// newUser.Password = string(newPassword)
			// fmt.Println(newUser.Password)
			newUser.Password = string(hashedPassword)
			response_Register := repository.UserRegisterRepository(h.db, newUser)
			jsonData, _ := json.Marshal(&response_Register)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(201)
			w.Write(jsonData)
		}
	}
}

type LoginHandler struct {
	db *sql.DB
}

// Login implements LoginHandlerIF
func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var newUser entity.User
		var validate *entity.User

		json.NewDecoder(r.Body).Decode(&newUser)
		fmt.Println(r.Body)
		tempPassword := newUser.Password
		newPassword := []byte(newUser.Password)
		_, err := bcrypt.GenerateFromPassword(newPassword, bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		validate = &newUser
		serv := service.NewUserService()

		// newUser.Password = string(hashedPassword)
		// fmt.Println(newUser.Password)
		sqlQuery := `select u.id, u.username, u.email, u.password, u.age,
				u.created_at, u.updated_at from public.users as u where email = $1`

		err = h.db.QueryRow(sqlQuery, newUser.Email).
			Scan(&newUser.Id, &newUser.Username, &newUser.Email, &newUser.Password,
				&newUser.Age, &newUser.CreatedAt, &newUser.UpdatedAt)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(newUser)
		validate, err = serv.Login(validate, tempPassword)
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

type UsersHandlerIF interface {
	UsersHandler(w http.ResponseWriter, r *http.Request)
}

type RegisterHandlerIF interface {
	Register(w http.ResponseWriter, r *http.Request)
}

type LoginHandlerIF interface {
	Login(w http.ResponseWriter, r *http.Request)
}

func NewUsersHandler(db *sql.DB) UsersHandlerIF {
	return &UsersHandler{db: db}
}

func NewRegisterHandler(db *sql.DB) RegisterHandlerIF {
	return &RegisterHandler{db: db}
}

func UserLoginHandler(db *sql.DB) LoginHandlerIF {
	return &LoginHandler{db: db}
}

func (h *UsersHandler) Updateusr(w http.ResponseWriter, r *http.Request, id string) {
	if id != "" {
		ctx := r.Context()
		user := middleware.RunUser(ctx)
		fmt.Println(user)
		fmt.Println(user.Id)
		var newUser entity.User
		json.NewDecoder(r.Body).Decode(&newUser)
		fmt.Println(newUser)
		var validasiUser *entity.User
		validasiUser = &newUser
		servic := service.NewUserService()
		validasiUser, err := servic.Update(validasiUser)
		if err != nil {
			fmt.Println(err)
		}

		newUser.UpdatedAt = time.Now()
		sqlQuery := `
		UPDATE public.users set email = $1, username= $2, updated_at = $3 
		where id = $4`

		res, err := h.db.Exec(sqlQuery,
			newUser.Email,
			newUser.Username,
			newUser.UpdatedAt,
			id,
		)
		fmt.Println(res)
		if err != nil {
			fmt.Println("error update")
			w.Write([]byte(fmt.Sprint(err)))

		}
		sqlQuery1 := `select u.id, u.username, u.email, u.password, u.age,
		u.created_at, u.updated_at from public.users as u  where id= $1`
		err = h.db.QueryRow(sqlQuery1, id).
			Scan(&newUser.Id, &newUser.Username, &newUser.Email, &newUser.Password,
				&newUser.Age, &newUser.CreatedAt, &newUser.UpdatedAt)
		// count, err := res.RowsAffected()
		if err != nil {
			w.Write([]byte(fmt.Sprint(err)))

		}

		fmt.Println(newUser)
		newUser.UpdatedAt = time.Now()
		responseUpdateUser := entity.ResponseUpdate{
			Id:        newUser.Id,
			Email:     newUser.Email,
			Username:  newUser.Username,
			Age:       newUser.Age,
			UpdatedAt: newUser.UpdatedAt,
		}
		jsonData, _ := json.Marshal(&responseUpdateUser)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(jsonData)
		return
	}
}

func (h *UsersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := middleware.RunUser(ctx)

	fmt.Println(user)
	fmt.Println(user.Id)
	// if temp_id != nil{}
	sqlstament := `DELETE from public.users where id = $1`
	_, err := h.db.Exec(sqlstament, user.Id)

	if err != nil {
		w.Write([]byte(fmt.Sprint(err)))

	}
	message := entity.Message{
		Message: "Your account has been successfully deleted",
	}
	jsonData, _ := json.Marshal(&message)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonData)

}
