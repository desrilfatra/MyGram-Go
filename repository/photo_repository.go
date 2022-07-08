package repository

import (
	"database/sql"
	"fmt"
	"log"
	"mygram-go/entity"
	"time"
)

func PhotoGetRepo(db *sql.DB) []*entity.ResponseGetPhoto {
	sqlQuery := `
	select p.id, p.title,p.caption, p.photo_url, p.user_id, p.created_at,
   	p.updated_at, u.email, u.username 
    from photo as p inner join users as u on p.user_id = u.id`
	rows, err := db.Query(sqlQuery)
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
	return photos
}

func PhotoPostRepo(db *sql.DB, newPhotos entity.Photo, user_id int) entity.ResponsePostPhoto {
	sqlQuery := `INSERT into photo
	(title,caption,photo_url,user_id,created_at,updated_at)
	values (?,?,?,?,?,?)`
	res, err := db.Exec(sqlQuery,
		newPhotos.Title,
		newPhotos.Caption,
		newPhotos.Url,
		user_id,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		panic(err.Error())
	}

	lastIdphoto, err := res.LastInsertId()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("The last inserted row id: %d\n", lastIdphoto)
	if err != nil {
		panic(err)
	}
	response := entity.ResponsePostPhoto{
		Id:        int(lastIdphoto),
		Title:     newPhotos.Title,
		Caption:   newPhotos.Caption,
		Url:       newPhotos.Url,
		User_id:   int(user_id),
		CreatedAt: time.Now(),
	}
	return response
}

func PhotoPutRepo(db *sql.DB, newPhotos entity.Photo, id string) entity.ResponsePutPhoto {
	sqlQuery := `update photo set title = ?, caption = ? , photo_url = ?, updated_at = ? where id = ?`
	_, err = db.Exec(sqlQuery,
		newPhotos.Title,
		newPhotos.Caption,
		newPhotos.Url,
		time.Now(),
		id,
	)
	if err != nil {
		fmt.Println("error update")
		panic(err)
	}
	sqlQuery1 := `select p.id, p.title, p.caption, p.photo_url, p.user_id, p.created_at,
				p.updated_at  from photo as p where p.id= ?`
	res, err := db.Query(sqlQuery1, id)
	defer res.Close()

	if err != nil {
		log.Fatal(err)
	}
	for res.Next() {
		err := res.Scan(&newPhotos.Id, &newPhotos.Title, &newPhotos.Caption, &newPhotos.Url,
			&newPhotos.User_id, &newPhotos.CreatedAt, &newPhotos.UpdatedAt)
		if err != nil {
			log.Fatal(err)
		}
	}
	response := entity.ResponsePutPhoto{
		Id:        newPhotos.Id,
		Title:     newPhotos.Title,
		Caption:   newPhotos.Caption,
		Url:       newPhotos.Url,
		User_id:   newPhotos.User_id,
		UpdatedAt: newPhotos.UpdatedAt,
	}
	return response
}

func PhotoDeleteRepo(db *sql.DB, id string) entity.Message {
	sqlQuery := `delete from public.photo where id = ?`
	_, err := db.Exec(sqlQuery, id)
	if err != nil {
		fmt.Println("error delete")
		panic(err)
	}
	return entity.Message{
		Message: "Your photo has been Successfully deleted",
	}
}
