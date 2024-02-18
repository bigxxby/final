package test

import (
	"log"
	"net/http"
	"test/internal/myDatabase"
	"text/template"
)

var db = myDatabase.NewDatabase("myDatabase.db")

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
	if r.Method == "GET" {

		data, err := template.ParseFiles("ui/static/templates/login.html")
		if err != nil {
			log.Println(err)
			return
		}
		data.Execute(w, nil)
	}
	if r.Method == "POST" {
		_, err := Login(w, r)
		if err != nil {
			log.Println(err.Error())

			w.Write([]byte("User not found or password incorrect"))
			return
		}
		data, err := template.ParseFiles("ui/static/templates/login.html")
		if err != nil {
			log.Println(err)
			return
		}
		data.Execute(w, nil)
	}
}

func RegHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/reg" {
		http.NotFound(w, r)
		log.Println("Not found")
		return
	}
	if r.Method == "GET" {

		data, err := template.ParseFiles("ui/static/templates/reg.html")
		if err != nil {
			log.Println(err)
			return
		}
		data.Execute(w, nil)
	}
	if r.Method == "POST" {
		_, err := Registration(w, r) // returning user
		if err != nil {
			// panic(err.Error())
			w.Write([]byte(err.Error()))
			return
		}
		// data, err := template.ParseFiles("ui/static/templates/reg.html")
		// if err != nil {
		// log.Println(err)
		// return
		// }
		// data.Execute(w, nil)
	}
}
