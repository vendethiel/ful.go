package main

/**
 * http user actions
 */

import (
	"net/http"
	"database/sql"

	"github.com/julienschmidt/httprouter"
)

type authHandler func(*User, http.ResponseWriter, *http.Request, httprouter.Params)

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
    user.Role = validateUserRole(user.Role, u)
		newUser, err := InsertUser(db, user)
		if err != nil {
			SendError(w, err)
			return
		}
		MarkAsJSON(w)
    w.WriteHeader(201)
		SendJSON(w, newUser)
	}
}

func UserUpdate(db *sql.DB) authHandler {
	return func(u *User, w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    MarkAsJSON(w)
		id := p.ByName("id")
		user, err := parseUser(r)
		if err != nil {
			SendError(w, err)
			return
		}
    if updateUser, err := GetUserBy(db, "id", id); err != nil || updateUser.Id == 0 {
      // no such user
      w.WriteHeader(404)
      return
    }
		if user.Role == "normal" || user.Role == "admin" {
			if u.Role == "admin" {
				UpdateUserColumn(db, id, "role", user.Role)
			} else {
        // non-admin can't update other users' roles
        w.WriteHeader(401)
        return
      }
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
    w.WriteHeader(204)
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
