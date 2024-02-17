package main

import (
	"fmt"
	"log"
	"net/http"
	test "test/internal"
)

func main() {
	// Создание или подключение к базе данных SQLite
	fmt.Println("Successfully connected to the database.")
	fs := http.FileServer(http.Dir("ui/static"))
	// Обработчик для всех запросов к статическим файлам
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Указываем маршрут для обслуживания статических файлов
	http.HandleFunc("/", test.MainHandler)
	http.HandleFunc("/login", test.LoginHandler)
	http.HandleFunc("/reg", test.RegHandler)

	log.Println("Server started on http://localhost:8080/")
	http.ListenAndServe(":8080", nil)

	// Здесь вы можете выполнить любые операции с вашей базой данных SQLite
	// Например, создать таблицы, вставить данные и т. д.
}
