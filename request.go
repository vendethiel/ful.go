package main

/**
 * request-parsing helpers
 */

import (
	"database/sql"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
)

func Authenticate(db *sql.DB, action authHandler) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		basicUser, errCode := parseAuthHeader(r.Header["Authorization"])
		if errCode != 200 {
			w.WriteHeader(int(errCode))
			return
		}

		user, err := GetUserBy(db, "email", basicUser.Username)
		if err != nil || user.Password != basicUser.Password {
			w.WriteHeader(401)
			return
		}

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
