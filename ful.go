package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)

func connectToDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:root@/ful")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	db, err := connectToDB()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/users/:id", UserShow(db))
	router.GET("/search/users", UserSearch(db))
	router.POST("/users", UserCreate(db))
	router.PUT("/users/:id", UserUpdate(db))
	router.DELETE("/users/:id", UserRemove(db))

	log.Print("Starting....\n")
	log.Fatal(http.ListenAndServe(":8080", router))
}
