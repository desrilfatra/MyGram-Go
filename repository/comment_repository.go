package repository

import (
	"database/sql"
	"fmt"
	"mygram-go/entity"
	"time"
)

func CommentGetRepository(db *sql.DB) []*entity.ResponseGetComment {
	sqlQuery := `
	select c.id, c.message,c.photo_id,c.user_id,c.updated_at,c.created_at,
	u.id,u.email,u.username,p.id,p.title,p.caption,p.photo_url,p.user_id 
	from comment c left join public.photo p on c.photo_id = p.id 
	left join users u on c.user_id = u.id`
	rows, err := db.Query(sqlQuery)
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
	return comments
}

func CommentPostRepository(db *sql.DB, comment entity.Commment, User_id int) entity.ResponsePostComment {
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	sqlQuery := `INSERT INTO public.comment
				(user_id,photo_id,message,created_at,updated_at)
				values ($1,$2,$3,$4,$5) Returning id`
	err := db.QueryRow(sqlQuery,
		User_id,
		comment.Photo_id,
		comment.Message,
		comment.CreatedAt,
		comment.UpdatedAt).Scan(&comment.Id)
	if err != nil {
		panic(err)
	}
	responseComment := entity.ResponsePostComment{
		Id:        comment.Id,
		Message:   comment.Message,
		Photo_id:  comment.Photo_id,
		User_id:   int(User_id),
		CreatedAt: time.Now(),
	}
	return responseComment
}

func CommentPutRepository(db *sql.DB, comment entity.Commment, id string) entity.ResponsePutComment {
	sqlQuery := `update public.comment set message = $1, updated_at =$2 where id = $3`
	//query.scan
	_, err = db.Exec(sqlQuery,
		comment.Message,
		time.Now(),
		id,
	)
	if err != nil {
		fmt.Println("error update")
		panic(err)
	}
	response := entity.ResponsePutComment{}
	sqlQuery1 := `SELECT c.id,p.title,p.caption,p.photo_url,c.user_id,c.updated_at 
	from comment c left join photo p on c.photo_id = p.id where c.id= $1`
	err = db.QueryRow(sqlQuery1, id).Scan(&response.Id, &response.Title,
		&response.Caption, &response.Url, &response.User_id, &response.UpdatedAt)
	if err != nil {
		panic(err)
	}
	return response
}

func CommentDeleteRepository(db *sql.DB, id string) entity.Message {
	sqlQuery := `DELETE from public.comment where id = $1`
	_, err := db.Exec(sqlQuery, id)
	if err != nil {
		panic(err)
	}
	response := entity.Message{
		Message: "Your Comment has been successfully deleted",
	}
	return response

}
