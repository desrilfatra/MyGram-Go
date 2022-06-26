package main

import (
	"database/sql"
	"fmt"
	"log"
	"mygram-go/controller"
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

	route.HandleFunc("/users", usersHandler.UsersHandler)
	route.HandleFunc("/users/register", registerHandler.RegisterUser)
	route.HandleFunc("/users/login", loginHandler.Login)
	route.HandleFunc("/users/{id}", usersHandler.UsersHandler)

	//handler photo
	photoHandler := controller.NewPhoto(db)
	route.HandleFunc("/photos", photoHandler.Photo)
	route.HandleFunc("/photos/{id}", photoHandler.Photo)

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
