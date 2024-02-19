package test

import (
	"net/http"
	"test/internal/myDatabase"
)

func Registration(w http.ResponseWriter, r *http.Request) (myDatabase.User, error) {
	user := myDatabase.User{}
	name := r.FormValue("name")
	surname := r.FormValue("surname")
	password := r.FormValue("password")
	email := r.FormValue("email")
	isAdminthis := 0
	if r.FormValue("checkbox") == "on" {
		isAdminthis = 1
	}
	user.Email = email
	user.Name = name
	user.Password = password
	user.Surname = surname
	user.IsAdmin = isAdminthis
	err := db.AddUser(name, surname, email, password, isAdminthis)
	if err != nil {
		return user, err
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return user, nil
}

