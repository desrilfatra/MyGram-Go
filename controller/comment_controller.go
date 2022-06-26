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
	serv := service.NewUserService()
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]
	temp_id := serv.VerivyToken(reqToken)
	fmt.Println(temp_id)
	switch r.Method {
	case http.MethodGet:
		fmt.Println("Get Comments")
		sqlQuery := `
		select c.id, c.caption,c.photo_id,c.user_id,c.updatedat,c.createdat,
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
			if scanerr := rows.Scan(&comment.Id, &comment.Message, &comment.Photo_id, &comment.User_id, &comment.UpdatedAt, &comment.CreatedAt, &comment.User.Id, &comment.User.Email, &comment.User.Username, &comment.Photo.Id, &comment.Photo.Title, &comment.Photo.Caption, &comment.Photo.Url, &comment.Photo.User_id); scanerr != nil {
				fmt.Println("Scan error", scanerr)
			}
			comments = append(comments, &comment)
		}
		jsonData, _ := json.Marshal(&comments)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(jsonData)

	}
}
