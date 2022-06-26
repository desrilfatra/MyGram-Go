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

type SosialMediaHandlerInterf interface {
	SosialMedia(w http.ResponseWriter, r *http.Request)
}

func NewSosialMedia(db *sql.DB) SosialMediaHandlerInterf {
	return &SosialMediaHand{db: db}
}

type SosialMediaHand struct {
	db *sql.DB
}

// SosialMedia implements SosialMediaHandlerInterf
func (sm *SosialMediaHand) SosialMedia(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	fmt.Println(id)
	serv := service.NewUserService()
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]
	temp_id := serv.VerivyToken(reqToken)
	fmt.Println(temp_id)
	switch r.Method {
	case http.MethodGet:
		fmt.Println("Get Social Media")
		sqlQuery := `
		select distinct on (sm.id) sm.id, sm.name, sm.sosial_media_url, sm.userid,
   		 u.createdat, u.updatedat, u.id, u.username, p.photo_url 
   		 from public.sosialmedia sm left join public.users u on sm.userid = u.id
   		 left join public.photo p on u.id = p.user_id  `
		rows, err := sm.db.Query(sqlQuery)
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()
		socialmedias := []*entity.ResponseGetSocialMedia{}
		for rows.Next() {
			var socialmedia entity.ResponseGetSocialMedia
			if serr := rows.Scan(&socialmedia.Id, &socialmedia.Name, &socialmedia.Social_Media_Url, &socialmedia.User_id,
				&socialmedia.CreatedAt, &socialmedia.UpdatedAt, &socialmedia.User.Id, &socialmedia.User.Username,
				&socialmedia.User.Url); serr != nil {
				fmt.Println("Scan error", serr)
			}
			socialmedias = append(socialmedias, &socialmedia)
		}
		jsonData, _ := json.Marshal(&socialmedias)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(jsonData)

	}
}
