package main

import (
	"database/sql"
	"fmt"
	"mygram-go/router"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB

	err error
)

func main() {
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/db-sql-go?parseTime=true")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Succesfully connected to database")
	router.RunRoute(db)

}
