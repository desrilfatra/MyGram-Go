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

type CommentHandlerInterf interface {
	Comment(w http.ResponseWriter, r *http.Request)
}

func NewComment(db *sql.DB) CommentHandlerInterf {
	return &CommentHand{db: db}
}

type CommentHand struct {
	db *sql.DB
}

// Comment implements CommentHandlerInterf
func (ch *CommentHand) Comment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	fmt.Println(id)
	ctx := r.Context()
	user := middleware.ForUser(ctx)

	fmt.Println(user)
	fmt.Println(user.Id)
	commentservic := service.CommentServic()
	switch r.Method {
	case http.MethodGet:
		fmt.Println("Get Comments")
		sqlQuery := `
		select c.id, c.caption,c.photo_id,c.user_id,c.updated_at,c.created_at,
   		u.id,u.email,u.username,p.id,p.title,p.caption,p.photo_url,p.user_id 
    	from comment c left join public.photo p on c.photo_id = p.id 
    	left join users u on c.user_id = u.id`
		rows, err := ch.db.Query(sqlQuery)
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()
		comments := []*entity.ResponseGetComment{}
		for rows.Next() {
			var comment entity.ResponseGetComment
			if scanerr := rows.Scan(&comment.Id, &comment.Message, &comment.Photo_id, &comment.User_id, &comment.UpdatedAt,
				&comment.CreatedAt, &comment.User.Id, &comment.User.Email, &comment.User.Username, &comment.Photo.Id,
				&comment.Photo.Title, &comment.Photo.Caption, &comment.Photo.Url, &comment.Photo.User_id); scanerr != nil {
				fmt.Println("Scan error", scanerr)
			}
			comments = append(comments, &comment)
		}
		jsonData, _ := json.Marshal(&comments)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(jsonData)

	case http.MethodPost:
		fmt.Println("POST")
		var newComment entity.Commment
		json.NewDecoder(r.Body).Decode(&newComment)
		err := commentservic.CekPostComment(newComment.Message)
		if err != nil {
			w.Write([]byte(fmt.Sprint(err)))
		} else {
			sqlStament := `Insert into public.comment
			(user_id,photo_id,message,created_at,updated_at)
			values ($1,$2,$3,$4,$4) Returning id`
			// intId, err := strconv.Atoi(id)
			err = ch.db.QueryRow(sqlStament,
				user.Id,
				newComment.Photo_id,
				newComment.Message,
				time.Now(),
			).Scan(&newComment.Id)
			if err != nil {
				w.Write([]byte(fmt.Sprint(err)))
			}
			response := entity.ResponsePostComment{
				Id:        newComment.Id,
				Message:   newComment.Message,
				Photo_id:  newComment.Photo_id,
				User_id:   int(user.Id),
				CreatedAt: time.Now(),
			}

			jsonData, _ := json.Marshal(&response)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(201)
			w.Write(jsonData)

		}
	}
}
