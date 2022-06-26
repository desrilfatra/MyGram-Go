package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"mygram-go/entity"
	"mygram-go/middleware"
	"mygram-go/service"
	"net/http"
	"time"

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
	ctx := r.Context()
	user := middleware.ForUser(ctx)

	fmt.Println(user)
	fmt.Println(user.Id)
	switch r.Method {
	case http.MethodGet:
		fmt.Println("Get Photo")
		sqlQuery := `
		select p.id, p.title,p.caption, p.photo_url, p.user_id, p.created_at,
   		p.updated_at, u.email, u.username 
    	from public.photo as p inner join public.users as u on p.user_id = u.id`
		rows, err := ph.db.Query(sqlQuery)
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()
		photos := []*entity.ResponseGetPhoto{}
		for rows.Next() {
			var photo entity.ResponseGetPhoto
			if scanerr := rows.Scan(&photo.Id, &photo.Title, &photo.Caption, &photo.Url, &photo.User_id,
				&photo.CreatedAt, &photo.UpdatedAt, &photo.Users.Email, &photo.Users.Username); scanerr != nil {
				fmt.Println("Scan error", scanerr)
			}
			photos = append(photos, &photo)
		}
		jsonData, _ := json.Marshal(&photos)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(jsonData)

	case http.MethodPost:
		fmt.Println("Post")

		var newPhotos entity.Photo
		json.NewDecoder(r.Body).Decode(&newPhotos)
		fmt.Println(newPhotos)
		//validasi user
		photoserv := service.NewPhotoService()
		err := photoserv.CekPostPhoto(newPhotos.Title, newPhotos.Url)
		if err != nil {
			w.Write([]byte(fmt.Sprint(err)))
		} else {
			//query insert
			sqlQuery1 := `insert into public.photo
			(title,caption,photo_url,user_id,created_at,updated_at)
			values ($1,$2,$3,$4,$5,$5) Returning id`
			//query.scan
			err = ph.db.QueryRow(sqlQuery1,
				newPhotos.Title,
				newPhotos.Caption,
				newPhotos.Url,
				user.Id,
				time.Now(),
			).Scan(&newPhotos.Id)
			if err != nil {
				w.Write([]byte(fmt.Sprint(err)))
			}
			response := entity.ResponsePostPhoto{
				Id:        newPhotos.Id,
				Title:     newPhotos.Title,
				Caption:   newPhotos.Caption,
				Url:       newPhotos.Url,
				User_id:   int(user.Id),
				CreatedAt: time.Now(),
			}

			jsonData, _ := json.Marshal(&response)
			w.Header().Add("Content-Type", "application/json")
			w.Write(jsonData)
			w.WriteHeader(201)
		}

	}
}

type PhotoHandlerInterf interface {
	Photo(w http.ResponseWriter, r *http.Request)
}

func NewPhoto(db *sql.DB) PhotoHandlerInterf {
	return &PhotoHand{db: db}
}
