package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type authHandler func(*User, http.ResponseWriter, *http.Request, httprouter.Params)

func Authenticate(db *sql.DB, action authHandler) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// validate we have the header...
		authHeader := r.Header["Authorization"]
		if len(authHeader) == 0 {
			w.WriteHeader(401)
			return
		}
		// validate the header's shape...
		auth := strings.SplitN(authHeader[0], " ", 2)
		if len(auth) != 2 || auth[0] != "Basic" {
			w.WriteHeader(400)
			return
		}

		// parse the header... validate the size
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 {
			w.WriteHeader(400)
			return
		}
		// validate the user
		user, err := GetUserBy(db, "email", pair[0])
		if err != nil || user.Password != pair[1] {
			w.WriteHeader(401)
			return
		}

		// call the action with the authentified user
		action(&user, w, r, p)
	}
}

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

func UserShow(db *sql.DB) authHandler {
	return func(user *User, w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
		MarkAsJSON(w)
		if user, err := GetUserBy(db, "id", p.ByName("id")); err != nil {
			w.WriteHeader(404)
		} else if user.Role == "admin" {
			w.WriteHeader(401)
			SendJSON(w, JsonError{401, "unauthorized"})
		} else {
			SendJSON(w, user)
		}
	}
}

func UserSearch(db *sql.DB) authHandler {
	return func(u *User, w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

func UserCreate(db *sql.DB) authHandler {
	return func(u *User, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		MarkAsJSON(w)
		user, err := parseUser(r)
		if err != nil {
			SendError(w, err)
			return
		}
		// only two roles allowed: normal and admin
		// also prevent non-admins from creating admin users
		if (user.Role != "normal" && user.Role != "admin") ||
			u.Role != "admin" {
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

func UserUpdate(db *sql.DB) authHandler {
	return func(u *User, w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
			// non-admin can't update other users' roles
			if u.Role == "admin" {
				UpdateUserColumn(db, id, "role", user.Role)
			}
		}
	}
}

func UserRemove(db *sql.DB) authHandler {
	return func(u *User, w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if err := DeleteUser(db, p.ByName("id")); err != nil {
			SendError(w, err)
			return
		}
	}
}
