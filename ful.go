package main

import (
	"log"
	"net/http"
	//"fmt"
	"database/sql"
	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)

type JsonError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func markAsJson(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func sendJson(w http.ResponseWriter, o interface{}) error {
	if json, err := json.Marshal(o); err == nil {
		w.Write(json)
		return nil
	} else {
		return err
	}
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(200)
}

func UserShow(db *sql.DB) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
		markAsJson(w)
		if user, err := GetUserBy(db, "id", p.ByName("id")); err != nil {
			w.WriteHeader(404)
		} else if user.Role == "admin" {
			w.WriteHeader(401)
			sendJson(w, JsonError{404, "not found"})
		} else {
			sendJson(w, user)
		}
	}
}

func UserSearch(db *sql.DB) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {

	}
}

func connectToDb() (*sql.DB, error) {
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
	db, err := connectToDb()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/users/:id", UserShow(db))
	router.GET("/search/users", UserSearch(db))

	log.Print("Starting....\n")
	log.Fatal(http.ListenAndServe(":8080", router))
}
