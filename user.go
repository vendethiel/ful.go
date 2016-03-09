package main

/**
 * DB-related query functions
 */

import (
	"database/sql"
	"fmt"
)

type User struct {
	Id        int    `json:"id"`
	Lastname  string `json:"lastname"`
	Firstname string `json:"firstname"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Role      string `json:"role"`
}

type NewUser struct {
	Id int64 `json:"id"`
}

func GetUserBy(db *sql.DB, column string, value string) (User, error) {
	var user User
	sql := "SELECT id, firstname, lastname, email, password, role FROM user WHERE " + column + " = ?"
	err := db.QueryRow(sql, value).Scan(&user.Id, &user.Firstname, &user.Lastname, &user.Email, &user.Password, &user.Role)
	return user, err
}

func InsertUser(db *sql.DB, user User) (*NewUser, error) {
	sql := `
    INSERT INTO user
    (firstname, lastname, email, password, role)
    VALUES (?, ?,             ?,        ?,    ?)
  `
	res, err := db.Exec(sql, user.Firstname, user.Lastname, user.Email, user.Password, user.Role)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	newUser := NewUser{id}
	return &newUser, nil
}

func UpdateUserColumn(db *sql.DB, id, column, value string) error {
	sql := fmt.Sprintf(`
    UPDATE user
    SET %s = ?
    WHERE id = ?
  `, column)
	_, err := db.Exec(sql, value, id)
	return err
}

func DeleteUser(db *sql.DB, id string) error {
	sql := "DELETE FROM user WHERE id = ?"
	_, err := db.Exec(sql, id)
	return err
}

func validateUserRole(role string, authUser *User) string {
  if authUser.Role != "admin" {
    return "normal"
  }
  if role == "normal" || role == "admin" {
    return role
  }
  return "normal"
}
