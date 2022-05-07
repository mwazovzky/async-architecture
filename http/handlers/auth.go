package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type User struct {
	ID       string `json:"id"`
	UUID     string `json:"uuid"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Token    string `json:"token"`
}

type AuthHandlers struct {
	db *sql.DB
}

const defaultRole = "employee"

func init() {
	rand.Seed(time.Now().UnixNano())
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

func (ah *AuthHandlers) Register(rw http.ResponseWriter, r *http.Request) {
	log.Println("Register Request")

	user := User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(rw, "Unable process request params", http.StatusBadRequest)
		return
	}

	user.UUID = generateUUID()
	user.Role = defaultRole

	err = createUser(ah.db, &user)
	log.Println(user)
	if err != nil {
		http.Error(rw, "Unable to create users", http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")

	e := json.NewEncoder(rw)
	err = e.Encode(user)
	if err != nil {
		http.Error(rw, "Unable to marshall json", http.StatusInternalServerError)
	}
}

func (ah *AuthHandlers) Login(rw http.ResponseWriter, r *http.Request) {
	log.Println("Login Request")

	req := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(rw, "Unable process request params", http.StatusBadRequest)
		return
	}

	user := User{}
	err = findUserByEmailAndPassword(ah.db, &user, req.Email, req.Password)
	if err != nil {
		http.Error(rw, "Access denied", http.StatusForbidden)
		return
	}

	user.Token = generateToken()
	err = setUserToken(ah.db, &user)
	if err != nil {
		http.Error(rw, "Unable to update user", http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")

	e := json.NewEncoder(rw)
	err = e.Encode(user)
	if err != nil {
		http.Error(rw, "Unable to marshall json", http.StatusInternalServerError)
	}
}

func (ah *AuthHandlers) Logout(rw http.ResponseWriter, r *http.Request) {
	log.Println("Logout Request")

	req := struct {
		Token string `json:"token"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(rw, "Unable process request params", http.StatusBadRequest)
		return
	}

	user := User{}
	err = findUserByToken(ah.db, &user, req.Token)
	if err != nil {
		http.Error(rw, "Unable to find user", http.StatusInternalServerError)
		return
	}

	err = deleteUserToken(ah.db, &user)
	if err != nil {
		http.Error(rw, "Unable to update user", http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")

	res := struct {
		Status string `json:"status"`
	}{
		Status: "success",
	}
	e := json.NewEncoder(rw)
	err = e.Encode(res)
	if err != nil {
		http.Error(rw, "Unable to marshall json", http.StatusInternalServerError)
	}
}

func (ah *AuthHandlers) Check(rw http.ResponseWriter, r *http.Request) {
	log.Println("Check Request")

	req := struct {
		Token string `json:"token"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(rw, "Access denied", http.StatusForbidden)
		return
	}

	user := User{}
	err = findUserByToken(ah.db, &user, req.Token)
	if err != nil {
		http.Error(rw, "Access denied", http.StatusForbidden)
		return
	}

	rw.Header().Add("Content-Type", "application/json")

	e := json.NewEncoder(rw)
	err = e.Encode(user)
	if err != nil {
		http.Error(rw, "Unable to marshall json", http.StatusInternalServerError)
	}
}

func (ah *AuthHandlers) UserUpdate(rw http.ResponseWriter, r *http.Request) {
	log.Println("UserUpdate Request")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable process request params", http.StatusBadRequest)
		return
	}

	req := struct {
		Role string `json:"role"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(rw, "Unable process request params", http.StatusBadRequest)
		return
	}

	err = updateUser(ah.db, id, req.Role)
	if err != nil {
		http.Error(rw, "Unable to update user", http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")

	res := struct {
		Status string `json:"status"`
	}{
		Status: "success",
	}
	e := json.NewEncoder(rw)
	err = e.Encode(res)
	if err != nil {
		http.Error(rw, "Unable to marshall json", http.StatusInternalServerError)
	}
}

func (ah *AuthHandlers) UserDelete(rw http.ResponseWriter, r *http.Request) {
	log.Println("UserUpdate Request")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable process request params", http.StatusBadRequest)
		return
	}

	err = deleteUser(ah.db, id)
	if err != nil {
		http.Error(rw, "Unable to update user", http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")

	res := struct {
		Status string `json:"status"`
	}{
		Status: "success",
	}
	e := json.NewEncoder(rw)
	err = e.Encode(res)
	if err != nil {
		http.Error(rw, "Unable to marshall json", http.StatusInternalServerError)
	}
}

func getUsers(db *sql.DB, users *[]User) error {
	var user User

	rows, err := db.Query("SELECT id, uuid, username, email, role FROM users")
	if err != nil {
		return err
	}

	for rows.Next() {
		err = rows.Scan(&user.ID, &user.UUID, &user.Username, &user.Email, &user.Role)
		if err != nil {
			return err
		}

		*users = append(*users, user)
	}

	return nil
}

func findUserByEmailAndPassword(db *sql.DB, user *User, email string, password string) error {
	var token sql.NullString
	sql := "SELECT id, uuid, username, email, role, token FROM users WHERE email=$1 AND password=$2;"
	row := db.QueryRow(sql, email, password)
	err := row.Scan(&user.ID, &user.UUID, &user.Username, &user.Email, &user.Role, &token)
	if err != nil {
		return err
	}

	if token.Valid {
		user.Token = token.String
	}

	return nil
}

func findUserByToken(db *sql.DB, user *User, token string) error {
	sql := "SELECT id, uuid, username, email, role, token FROM users WHERE token=$1;"
	row := db.QueryRow(sql, token)
	err := row.Scan(&user.ID, &user.UUID, &user.Username, &user.Email, &user.Role, &user.Token)
	if err != nil {
		return err
	}

	return nil
}

func createUser(db *sql.DB, user *User) error {
	sql := "INSERT INTO users (uuid, username, email, password, role) VALUES($1,$2,$3,$4,$5) RETURNING id"
	err := db.QueryRow(sql, user.UUID, user.Username, user.Email, user.Password, user.Role).Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}

func updateUser(db *sql.DB, id int, role string) error {
	sql := "UPDATE users SET role=$1 WHERE id=$2;"
	err := db.QueryRow(sql, role, id).Err()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func setUserToken(db *sql.DB, user *User) error {
	sql := "UPDATE users SET token=$1 WHERE id=$2;"
	err := db.QueryRow(sql, user.Token, user.ID).Err()
	if err != nil {
		return err
	}

	return nil
}

func deleteUserToken(db *sql.DB, user *User) error {
	var token sql.NullString
	sql := "UPDATE users SET token=$1 WHERE id=$2;"
	err := db.QueryRow(sql, token, user.ID).Err()
	if err != nil {
		return err
	}

	return nil
}

func deleteUser(db *sql.DB, id int) error {
	sql := "DELETE FROM users WHERE id=$1;"
	err := db.QueryRow(sql, id).Err()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func generateUUID() string {
	return uuid.New().String()
}

func generateToken() string {
	charset := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	n := 24
	token := make([]rune, n)
	for i := range token {
		token[i] = charset[rand.Intn(len(charset))]
	}
	return string(token)
}
