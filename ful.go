package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
  "log"
  "fmt"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(200)
}

func UserShow(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
  fmt.Fprintf(w, "Yo %s\n", p.ByName("id"))
}

func main() {
  router := httprouter.New()
  router.GET("/", Index)
  router.GET("/users/:id", UserShow)

	log.Fatal(http.ListenAndServe(":8080", router))
}
