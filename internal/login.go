package test

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"test/internal/myDatabase"

	"github.com/gofrs/uuid"
)

func Login(w http.ResponseWriter, r *http.Request) (*myDatabase.User, error) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := db.Authenticate(email, password)
	if err != nil {
		log.Println("Error logging user: ", err.Error())
		return nil, err
	}
	sessionid, err := generateUuid()
	if err != nil {
		log.Println("error generating uuid")
		return nil, err
	}
	err = db.UpdateUserSession(int64(user.ID), sessionid)
	if err != nil {
		log.Println("error updating uuid")
		return nil, err
	}
	cookie := http.Cookie{
		Name:  "session_id",
		Value: sessionid,
	}
	http.SetCookie(w, &cookie)
	user.IsLogged = true
	user.SessionId = sql.NullString{String: sessionid}
	if err != nil {
		log.Println(err.Error())

		return nil, errors.New("Internal server error")
	}

	return user, nil
}

func generateUuid() (string, error) {
	sessionId, err := uuid.NewV4()
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	return sessionId.String(), nil
}

func isAuthenticated(r *http.Request) (bool, *myDatabase.User) {
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		fmt.Println("COOCKIE ERROR")
		return false, nil
	}
	userID := db.FindUserIdBySessionId(cookie.Value)
	if userID == -1 {
		log.Println("NOT FOUND COOCKIES ERROR")
		// Сессия не найдена в хранилище, пользователь не аутентифицирован.
		return false, nil
	}
	// Получаем информацию о пользователе из вашего хранилища.
	user, err := db.GetUserByID(userID)
	if err != nil {
		log.Println("NOT FOUND USER ERROR")
		// Ошибка при получении информации о пользователе.
		return false, nil
	}

	log.Println("ALL GOOD 200")
	// Пользователь аутентифицирован.
	user.Password = ""
	user.IsLogged = true
	return true, user
}
