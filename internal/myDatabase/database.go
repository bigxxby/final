package myDatabase

import (
	"database/sql"
	"log"
)
import _ "github.com/mattn/go-sqlite3"

type User struct {
	ID       int
	Name     string
	Surname  string
	Email    string
	Password string
}
type Database struct {
	Connection *sql.DB
}

func (db *Database) Close() {
	db.Connection.Close()
}
func NewDatabase(dbname string) *Database {

	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	varDb := &Database{Connection: db}
	err = varDb.InitializeUserTable()
	if err != nil {
		panic(err.Error())
		return nil
	}
	return varDb
}
func (db *Database) AddUser(name, surname, email, password string) error {
	// Подготовленный SQL-запрос для вставки данных
	query := "INSERT INTO users (name, surname, email, password) VALUES (?, ?, ?, ?)"

	// Выполнение запроса с передачей параметров
	_, err := db.Connection.Exec(query, name, surname, email, password)
	if err != nil {
		return err
	}
	log.Println("Added user:", name)
	return nil
}
func (db *Database) InitializeUserTable() error {
	// SQL запрос для создания таблицы пользователей, если она не существует
	query := `
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name VARCHAR(50) NOT NULL,
            surname VARCHAR(50) NOT NULL,
            email VARCHAR(100) NOT NULL UNIQUE,
            password VARCHAR(100) NOT NULL
        )
    `

	// Выполнение SQL запроса
	_, err := db.Connection.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetUserByID(userID int) (*User, error) {
	query := "SELECT id, name, surname, email, password FROM users WHERE id = ?"
	row := db.Connection.QueryRow(query, userID)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
