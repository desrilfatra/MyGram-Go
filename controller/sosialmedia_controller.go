package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"mygram-go/entity"
	"mygram-go/middleware"
	"mygram-go/service"
	"net/http"

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

	ctx := r.Context()
	user := middleware.ForUser(ctx)

	fmt.Println(user)
	fmt.Println(user.Id)
	sosialmediaser := service.SocialMediaSerN()

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

	case http.MethodPost:
		fmt.Println("POST")
		var newSocialMedia entity.SocialMedia
		json.NewDecoder(r.Body).Decode(&newSocialMedia)
		err := sosialmediaser.CekPostSocialMedia(newSocialMedia.Name, newSocialMedia.Social_Media_Url)
		if err != nil {
			w.Write([]byte(fmt.Sprint(err)))
		} else {
			sqlQuery := `insert into public.socialmedia
			(name,social_media_url,userid)
			values ($1,$2,$3) Returning id`
			// intId, err := strconv.Atoi(id)
			err = sm.db.QueryRow(sqlQuery, newSocialMedia.Name, newSocialMedia.Social_Media_Url, user.Id).Scan(&newSocialMedia.Id)
			if err != nil {
				fmt.Println(err)
			}

			response := entity.ResponsePostSocialMedia{}
			sqlQuery1 := `select s.id,s.name,s.social_media_url,s.userid,u.created_at 
			from public.socialmedia s left join public.users u on s.userid = u.id where s.id = $1`
			err = sm.db.QueryRow(sqlQuery1, newSocialMedia.Id).Scan(&response.Id, &response.Name,
				&response.Social_Media_Url, &response.User_id, &response.CreatedAt)
			if err != nil {
				w.Write([]byte(fmt.Sprint(err)))
			}

			jsonData, _ := json.Marshal(&response)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(201)
			w.Write(jsonData)
		}

	case http.MethodPut:
		fmt.Println("PUT")
		if id != "" {
			var newSocialMedia entity.SocialMedia
			json.NewDecoder(r.Body).Decode(&newSocialMedia)
			err := sosialmediaser.CekPostSocialMedia(newSocialMedia.Name, newSocialMedia.Social_Media_Url)
			if err != nil {
				w.Write([]byte(fmt.Sprint(err)))
			} else {
				sqlQuery := `update public.socialmedia set name = $1, social_media_url = $2 where id = $3`
				//query.scan
				_, err = sm.db.Exec(sqlQuery,
					newSocialMedia.Name,
					newSocialMedia.Social_Media_Url,
					id,
				)
				if err != nil {
					fmt.Println("error update")
					w.Write([]byte(fmt.Sprint(err)))
				}
				response := entity.ResponsePutSocialMedia{}
				sqlQuery1 := `select s.id,s.name, s.social_media_url, s.userid, u.updated_at 
				from public.socialmedia s left join public.users u on s.userid = u.id where s.id = $1`
				err = sm.db.QueryRow(sqlQuery1, id).
					Scan(&response.Id, &response.Name, &response.Social_Media_Url, &response.User_id, &response.UpdatedAt)
				// count, err := res.RowsAffected()
				if err != nil {
					w.Write([]byte(fmt.Sprint(err)))
				}
				jsonData, _ := json.Marshal(&response)
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(200)
				w.Write(jsonData)
			}

		}

	}

}
