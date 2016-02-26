package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
  "log"
  "fmt"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)

type User struct {
  id int
  firstname, lastname, email, password, role string
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(200)
}

func UserShow(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
  err := db.QueryRow("SELECT firstname, lastname, email, password, role FROM user WHERE id = ?", p.ByName("id"))
}

func connectToDb() (*sql.DB, error) {
  db, err := sql.Open("mysql", "root:root@/ful")
  if err != nil {
    return nil, err
  }
  defer db.Close()

  err = db.Ping()
  if err != nil {
    return nil, err
  }
  return db, nil
}

func main() {
  _, err := connectToDb()
  if err != nil {
    panic(err.Error())
  }

  router := httprouter.New()
  router.GET("/", Index)
  router.GET("/users/:id", UserShow)

	log.Fatal(http.ListenAndServe(":8080", router))
}
