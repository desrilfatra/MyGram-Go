package repository

import (
	"database/sql"
	"fmt"
	"mygram-go/entity"
)

func SocmedGetRepo(db *sql.DB) []*entity.ResponseGetSocialMedia {
	sqlQuery := `
		select sm.id, sm.name, sm.social_media_url, sm.userid,
   		u.created_at, u.updated_at, u.id, u.username, p.photo_url 
   		from socialmedia as sm left join users as u on sm.userid = u.id
   		left join photo p on u.id = p.user_id `
	rows, err := db.Query(sqlQuery)
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
	return socialmedias
}

func SocmedPostRepo(db *sql.DB, newSocialMedia entity.SocialMedia, user_id int) entity.ResponsePostSocialMedia {
	sqlQuery := `INSERT into socialmedia (name,social_media_url,userid) values (?,?,?)`
	// intId, err := strconv.Atoi(id)
	res, err := db.Exec(sqlQuery,
		newSocialMedia.Name,
		newSocialMedia.Social_Media_Url,
		user_id)
	if err != nil {
		panic(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}
	response := entity.ResponsePostSocialMedia{}
	sqlQuery1 := `select s.id,s.name, s.social_media_url, s.userid
	from socialmedia s left join users u on s.userid = u.id where s.id = ?`
	err = db.QueryRow(sqlQuery1, id).
		Scan(&response.Id, &response.Name, &response.Social_Media_Url, &response.User_id)
	// count, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	return response
}

func SocmedPutRepo(db *sql.DB, SocialMediap entity.SocialMedia, id string) entity.ResponsePutSocialMedia {
	sqlQuery := `update socialmedia set name = ?, social_media_url = ? where id = ?`
	//query.scan
	_, err = db.Exec(sqlQuery,
		SocialMediap.Name,
		SocialMediap.Social_Media_Url,
		id,
	)
	if err != nil {
		fmt.Println("error update")
		panic(err)
	}
	response := entity.ResponsePutSocialMedia{}
	sqlQuery1 := `select s.id,s.name, s.social_media_url, s.userid, u.updated_at 
	from socialmedia s left join users u on s.userid = u.id where s.id = ?`
	err = db.QueryRow(sqlQuery1, id).Scan(&response.Id, &response.Name, &response.Social_Media_Url, &response.User_id, &response.UpdatedAt)
	// count, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	return response
}

func SocmedDelRepo(db *sql.DB, Id string) entity.Message {
	sqlQuery := `DELETE from socialmedia where id = ?`
	_, err := db.Exec(sqlQuery, Id)

	if err != nil {
		panic(err)
	}
	response := entity.Message{
		Message: "Your socialmedia has been successfully deleted",
	}
	return response
}
