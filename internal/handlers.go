package test

import (
	"log"
	"net/http"
	"text/template"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		log.Println("Not found")

		return
	}
	if r.Method != "GET" {
		w.Write([]byte("Method not allowed"))
		return
	}
	data, err := template.ParseFiles("ui/static/templates/main.html")
	if err != nil {
		log.Println(err)
		return
	}
	data.Execute(w, nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		http.NotFound(w, r)
		log.Println("Not found")

		return
	}
	data, err := template.ParseFiles("ui/static/templates/login.html")
	if err != nil {
		log.Println(err)
		return
	}
	data.Execute(w, nil)
}

func RegHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/reg" {
		http.NotFound(w, r)
		log.Println("Not found")

		return
	}
	data, err := template.ParseFiles("ui/static/templates/reg.html")
	if err != nil {
		log.Println(err)
		return
	}
	data.Execute(w, nil)
}
