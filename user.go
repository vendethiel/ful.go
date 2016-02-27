package main

import (
  "database/sql"
)

type User struct {
  Id int `json:"id"`
  Lastname string `json:"lastname"`
  Firstname string `json:"firstname"`
  Email string `json:"email"`
  Password string `json:"-"`
  Role string `json:"role"`
}

func GetUserBy(db *sql.DB, column string, strid string) (User, error) {
  var (
    id int
    firstname, lastname, email, password, role string
  )
  sql := "SELECT id, firstname, lastname, email, password, role FROM user WHERE " + column + " = ?"
  err := db.QueryRow(sql, strid).Scan(&id, &firstname, &lastname, &email, &password, &role)
  if err != nil {
    return User{}, err
  }
  return User{id, firstname, lastname, email, password, role}, nil
}
