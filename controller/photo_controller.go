package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"mygram-go/entity"
	"mygram-go/service"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type PhotoHand struct {
	db *sql.DB
}

// Photo implements PhotoHandlerInterf
func (ph *PhotoHand) Photo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	fmt.Println(id)

	switch r.Method {
	case http.MethodGet:
		fmt.Println("Get")
		servic := service.NewUserService()
		reqToken := r.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")
		reqToken = splitToken[1]
		temp_id := servic.VerivyToken(reqToken)
		fmt.Println(temp_id)
		sqlStament := `
		select p.id, p.title,p.caption, p.photo_url, p.user_id, p.createdat,
   		p.updatedat, u.email, u.username 
    	from public.photo as p inner join public.users as u on p.user_id = u.id`
		rows, err := ph.db.Query(sqlStament)
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()
		photos := []*entity.ResponseGetPhoto{}
		for rows.Next() {
			var photo entity.ResponseGetPhoto
			if serr := rows.Scan(&photo.Id, &photo.Title, &photo.Caption, &photo.Url, &photo.User_id, &photo.CreatedAt, &photo.UpdatedAt, &photo.Users.Email, &photo.Users.Username); serr != nil {
				fmt.Println("Scan error", serr)
			}
			photos = append(photos, &photo)
		}
		jsonData, _ := json.Marshal(&photos)
		w.Header().Add("Content-Type", "application/json")
		w.Write(jsonData)
		w.WriteHeader(200)

	}
}

type PhotoHandlerInterf interface {
	Photo(w http.ResponseWriter, r *http.Request)
}

func NewPhoto(db *sql.DB) PhotoHandlerInterf {
	return &PhotoHand{db: db}
}
