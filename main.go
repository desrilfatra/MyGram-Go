package main

import (
	"database/sql"
	"fmt"
	"log"
	"mygram-go/controller"
	"mygram-go/middleware"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "db-sql-go"
)

var (
	db *sql.DB

	err error
)

const PORT = ":8080"

func main() {
	db, err = sql.Open("postgres", ConnectDb(host, user, password, dbname, port))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Succesfully connected to database")

	// handler user
	route := mux.NewRouter()
	usersHandler := controller.NewUsersHandler(db)
	registerHandler := controller.NewRegisterHandler(db)
	loginHandler := controller.UserLoginHandler(db)

	route.HandleFunc("/users/register", registerHandler.Register)
	route.HandleFunc("/users/login", loginHandler.Login)
	route.Handle("/users", middleware.Auth(http.HandlerFunc(usersHandler.UsersHandler))).Methods("PUT")
	route.Handle("/users/{id}", middleware.Auth(http.HandlerFunc(usersHandler.UsersHandler))).Methods("Delete")

	//handler photo
	photoHandler := controller.NewPhoto(db)
	route.Handle("/photos", middleware.Auth(http.HandlerFunc(photoHandler.Photo))).Methods("GET")
	route.Handle("/photos", middleware.Auth(http.HandlerFunc(photoHandler.Photo))).Methods("POST")
	route.Handle("/photos/{id}", middleware.Auth(http.HandlerFunc(photoHandler.Photo))).Methods("PUT")
	route.Handle("/photos/{id}", middleware.Auth(http.HandlerFunc(photoHandler.Photo))).Methods("DELETE")

	//handler comment
	commentHandler := controller.NewComment(db)
	route.Handle("/comments", middleware.Auth(http.HandlerFunc(commentHandler.Comment))).Methods("GET")
	route.Handle("/comments", middleware.Auth(http.HandlerFunc(commentHandler.Comment))).Methods("POST")
	route.Handle("/comments/{id}", middleware.Auth(http.HandlerFunc(commentHandler.Comment))).Methods("PUT")

	//handler comment
	sosialmediaHandler := controller.NewSosialMedia(db)
	route.Handle("/sosialmedias", middleware.Auth(http.HandlerFunc(sosialmediaHandler.SosialMedia))).Methods("GET")
	route.Handle("/sosialmedias", middleware.Auth(http.HandlerFunc(sosialmediaHandler.SosialMedia))).Methods("POST")
	route.Handle("/sosialmedias/{id}", middleware.Auth(http.HandlerFunc(sosialmediaHandler.SosialMedia))).Methods("PUT")
	fmt.Println("Now listening on port 0.0.0.0" + PORT)
	srv := &http.Server{
		Handler:      route,
		Addr:         "0.0.0.0" + PORT,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
func ConnectDb(host, user, password, name string, port int) string {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbname)
	return psqlInfo
}
