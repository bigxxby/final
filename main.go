package main

import (
	"fmt"
	"net/http"
	"text/template"
)

func main() {
	staticDir := "./static"

	// Настройка обработчика для статических файлов
	fs := http.FileServer(http.Dir(staticDir))

	// Указываем маршрут для обслуживания статических файлов
	http.Handle("/", fs)

	// Обработчик для отправки пользовательских заголовков в ответ

	// Обработчик для вывода всех заголовков запроса HTTP
	http.HandleFunc("/requestheaders", requestHeadersHandler)

	// Обработчик для вывода HTTP метода, использованного в запросе
	http.HandleFunc("/httpmethod", httpMethodHandler)

	// Обработчик для чтения файла в формате JSON и отправки объекта с заданным ID
	http.HandleFunc("/userdata", userDataHandler)

	// Обработчик для отображения всех названий и значений заголовков в запросе
	http.HandleFunc("/allheaders", allHeadersHandler)
	http.HandleFunc("/main", mainHandler)

	// Запуск HTTP сервера на порту 8000
	http.ListenAndServe(":8080", nil)
}
func mainHandler(w http.ResponseWriter, r *http.Request) {
	file, _ := template.ParseFiles("static/index.html")
	err := file.Execute(w, nil)
	if err != nil {
		return
	}
}

func processDataHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем параметры из строки запроса
	name := r.URL.Query().Get("name")
	age := r.URL.Query().Get("age")
	city := r.URL.Query().Get("city")

	// Записываем данные в тело ответа HTTP
	fmt.Fprintf(w, "Name: %s, Age: %s, City: %s", name, age, city)
}

func customHeadersHandler(w http.ResponseWriter, r *http.Request) {
	// Отправляем пользовательские заголовки в ответ
	w.Header().Set("Content-Length", "123")
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "Custom headers sent successfully!")
}

func requestHeadersHandler(w http.ResponseWriter, r *http.Request) {
	// Выводим все заголовки запроса HTTP
	for name, headers := range r.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func httpMethodHandler(w http.ResponseWriter, r *http.Request) {
	// Выводим HTTP метод, использованный в запросе
	fmt.Fprintf(w, "HTTP Method: %s", r.Method)
}

func userDataHandler(w http.ResponseWriter, r *http.Request) {
	// Реализация этой функции зависит от вашей конкретной задачи чтения файла и отправки данных
	// Например, можно прочитать файл с данными пользователей, найти объект с заданным ID и отправить его
}

func allHeadersHandler(w http.ResponseWriter, r *http.Request) {
	// Выводим все заголовки запроса HTTP
	for name, headers := range r.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}
