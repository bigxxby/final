package test

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"test/internal/myDatabase"
	"time"

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
		Name:    "session_id",
		Value:   sessionid,
		Expires: time.Now().Add(24 * time.Hour),
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
