package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func parseUser(r *http.Request) (User, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return User{}, err
	}
	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		return User{}, err
	}
	return user, err
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}

func UserShow(db *sql.DB) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
		MarkAsJSON(w)
		if user, err := GetUserBy(db, "id", p.ByName("id")); err != nil {
			w.WriteHeader(404)
		} else if user.Role == "admin" {
			w.WriteHeader(401)
			SendJSON(w, JsonError{404, "not found"})
		} else {
			SendJSON(w, user)
		}
	}
}

func UserSearch(db *sql.DB) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		MarkAsJSON(w)
		email := r.URL.Query().Get("q")
		if len(email) == 0 {
			w.WriteHeader(404)
		} else if user, err := GetUserBy(db, "email", email); err != nil {
			w.WriteHeader(404)
		} else {
			SendJSON(w, []User{user})
		}
	}
}

func UserCreate(db *sql.DB) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		MarkAsJSON(w)
		user, err := parseUser(r)
		if err != nil {
			SendError(w, err)
			return
		}
		if user.Role != "normal" && user.Role != "admin" {
			user.Role = "normal"
		}
		newUser, err := InsertUser(db, user)
		if err != nil {
			SendError(w, err)
			return
		}
		MarkAsJSON(w)
		SendJSON(w, newUser)
	}
}

func UserUpdate(db *sql.DB) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id := p.ByName("id")
		user, err := parseUser(r)
		if err != nil {
			SendError(w, err)
			return
    }
    if user.Lastname != "" {
      UpdateUserColumn(db, id, "lastname", user.Lastname)
    }
    if user.Firstname != "" {
      UpdateUserColumn(db, id, "firstname", user.Firstname)
    }
    if user.Email != "" {
      UpdateUserColumn(db, id, "email", user.Email)
    }
    if user.Password != "" {
      UpdateUserColumn(db, id, "password", user.Password)
    }
    if user.Role == "normal" || user.Role == "admin" {
      UpdateUserColumn(db, id, "role", user.Role)
    }
	}
}

func UserRemove(db *sql.DB) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    if err := DeleteUser(db, p.ByName("id")); err != nil {
      SendError(w, err)
      return
    }
  }
}
