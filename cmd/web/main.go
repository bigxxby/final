package main

import (
	"fmt"
	"log"
	"net/http"

	test "test/internal"
)

func main() {
	fmt.Println("Successfully connected to the database.")
	fs := http.FileServer(http.Dir("ui/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", test.MainHandler)
	http.HandleFunc("/login", test.LoginHandler)
	http.HandleFunc("/reg", test.RegHandler)
	http.HandleFunc("/admin", test.AdminPage)
	http.HandleFunc("/delete", test.DeleteHandler)
	http.HandleFunc("/create", test.CreateUser)
	http.HandleFunc("/update", test.UpdateUser)
	http.HandleFunc("/message", test.MessageHandler)
	http.HandleFunc("/chord", test.ChordHandler)
	http.HandleFunc("/guitar", test.GuitarHandler)
	http.HandleFunc("/beginner", test.BeginnerHandler)
	http.HandleFunc("/metronome", test.MetronomeHandler)
	log.Println("Server started on http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
