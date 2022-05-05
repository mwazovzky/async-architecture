package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type User struct {
	UUID     string
	Username string
	Email    string
	Role     string
}

func init() {
	godotenv.Load()
}

func main() {
	db := connectDB()
	defer db.Close()

	users := []User{}
	err := getUsers(db, &users)
	if err != nil {
		log.Println(err)
	}

	log.Println(users)
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

func connectDB() *sql.DB {
	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")
	dbuser := os.Getenv("DB_USER")
	dbpassword := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbhost, dbport, dbuser, dbpassword, dbname)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to db")

	return db
}
