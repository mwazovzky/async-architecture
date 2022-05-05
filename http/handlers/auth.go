package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type User struct {
	UUID     string
	Username string
	Email    string
	Role     string
}

type AuthHandlers struct {
	db *sql.DB
}

func NewAuthHandlers(db *sql.DB) *AuthHandlers {
	return &AuthHandlers{db}
}

func (ah *AuthHandlers) UserIndex(rw http.ResponseWriter, r *http.Request) {
	log.Println("UserIndex Request")

	users := []User{}
	err := getUsers(ah.db, &users)
	if err != nil {
		http.Error(rw, "Unable to get users", http.StatusInternalServerError)
	}

	rw.Header().Add("Content-Type", "application/json")

	e := json.NewEncoder(rw)
	err = e.Encode(users)

	if err != nil {
		http.Error(rw, "Unable to marshall json", http.StatusInternalServerError)
	}
}

func getUsers(db *sql.DB, users *[]User) error {
	var user User

	rows, err := db.Query("SELECT uuid, username, email, role FROM users")
	if err != nil {
		return err
	}

	for rows.Next() {
		err = rows.Scan(&user.UUID, &user.Username, &user.Email, &user.Role)
		if err != nil {
			return err
		}

		*users = append(*users, user)
	}

	return nil
}
