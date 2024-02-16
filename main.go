package main

import (
	"net/http"
	"text/template"
)

func main() {
	staticDir := "./static"

	// Настройка обработчика для статических файлов
	fs := http.FileServer(http.Dir(staticDir))

	// Указываем маршрут для обслуживания статических файлов
	http.Handle("/", fs)

	http.ListenAndServe(":8080", nil)
}
func mainHandler(w http.ResponseWriter, r *http.Request) {
	file, _ := template.ParseFiles("static/index.html")
	err := file.Execute(w, nil)
	if err != nil {
		return
	}
}
func regHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		file, _ := template.ParseFiles("static/reg.html")
		err := file.Execute(w, nil)
		if err != nil {
			return
		}
	} else if r.Method == "POST" {
		firstName := r.FormValue("first name")
		lastName := r.FormValue("last name")
		email := r.FormValue("email")
		age := r.FormValue("age")
		password := r.FormValue("password")

		file, _ := template.ParseFiles("static/login.html")
		err := file.Execute(w, nil)
		if err != nil {
			return
		}
	}
}
