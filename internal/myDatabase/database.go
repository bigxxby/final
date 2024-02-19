package myDatabase

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID        int
	Name      string
	Surname   string
	Email     string
	Password  string
	IsAdmin   int
	IsLogged  bool
	SessionId sql.NullString
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
	}
	log.Println("Tablase users Created")
	return varDb
}

func (db *Database) AddUser(name, surname, email, password string, isAdmin int) error {
	query := "INSERT INTO users (name, surname, email, password, is_admin ) VALUES (?, ?, ?, ?, ? )"
	_, err := db.Connection.Exec(query, name, surname, email, password, isAdmin)
	if err != nil {
		log.Println(err, "//////////////////////")
		return err
	}
	log.Println("Added user:", name, surname, email, password, isAdmin)
	return nil
}

func (db *Database) InitializeUserTable() error {
	// SQL запрос для создания таблицы пользователей, если она не существует
	query := `
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name VARCHAR(50) NOT NULL,
            surname VARCHAR(50) NOT NULL,
            email VARCHAR(100) NOT NULL ,
            password VARCHAR(100) NOT NULL,
			is_admin INTEGER,
			session_id TEXT
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
	query := "SELECT * FROM users WHERE id = ?"
	row := db.Connection.QueryRow(query, userID)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &user.IsAdmin, &user.SessionId)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *Database) Authenticate(email string, password string) (*User, error) {
	query := "SELECT * FROM users WHERE email = ? AND password = ?"
	row := db.Connection.QueryRow(query, email, password)
	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &user.IsAdmin, &user.SessionId)
	if err != nil {
		return nil, err
	}
	user.Password = "" // Не рекомендуется возвращать пароль пользователя в результате аутентификации
	return &user, nil
}

func (db *Database) UpdateUserSession(userID int64, sessionID string) error {
	stmt, err := db.Connection.Prepare("UPDATE users SET session_id = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("failed to prepare SQL statement: %v", err)
	}
	defer stmt.Close()

	// Выполнение SQL-запроса для обновления Session ID пользователя
	_, err = stmt.Exec(sessionID, userID)
	if err != nil {
		return fmt.Errorf("failed to update user session ID: %v", err)
	}

	fmt.Printf("Session ID updated for user with ID %d\n", userID)
	return nil
}

func (db *Database) FindUserIdBySessionId(sessionId string) int {
	var id int
	// rows, err := db.Connection.Query("SELECT * FROM users WHERE session_id = ?", sessionId)
	err := db.Connection.QueryRow("SELECT id FROM users WHERE session_id = ?", sessionId).Scan(&id)
	if err != nil {
		log.Println(err.Error())

		return -1
	}
	return id
}

func (db *Database) GetAllUsers() *[]User {
	rows, err := db.Connection.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &user.IsAdmin, &user.SessionId)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	// Обработка ошибок после завершения итерации
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return &users
}

func (db *Database) DeleteUserById(userId string) {
	stmt, err := db.Connection.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Выполняем запрос DELETE
	_, err = stmt.Exec(userId)
	if err != nil {
		log.Fatal(err)
	}
}
