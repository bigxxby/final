package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"test/internal/myDatabase"
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
	_, user := isAuthenticated(r)

	data, err := template.ParseFiles("ui/static/templates/main.html")
	if err != nil {
		log.Println(err)
		return
	}
	data.Execute(w, user)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		http.NotFound(w, r)
		log.Println("Not found")

		return
	}
	if r.Method == "GET" {
		// db.FindUserBySessionId()

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

func AdminPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/admin" {
		http.NotFound(w, r)
		log.Println("Not found")

		return
	}
	if r.Method != "GET" {
		w.Write([]byte("Method not allowed"))
		return
	}
	boolean, user := isAuthenticated(r)
	if user == nil {
		http.Error(w, "Unathorised", http.StatusUnauthorized)
		return
	}
	if user.IsAdmin == 0 || !boolean {
		http.Error(w, "Unathorised", http.StatusUnauthorized)
		return
	}
	users := db.GetAllUsers()
	log.Println(users)

	data, err := template.ParseFiles("ui/static/templates/adminPanel.html")
	if err != nil {
		log.Println(err)
		return
	}
	data.Execute(w, users)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		id := r.Form.Get("id")
		fmt.Println("ID to delete:", id)
		db.DeleteUserById(id)
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
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
		log.Println(err.Error())
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	name := r.FormValue("name")
	surname := r.FormValue("surname")
	email := r.FormValue("email")
	password := r.FormValue("password")
	isAdmin := 0
	//  := r.FormValue("password")
	// password := r.FormValue("password")
	if r.FormValue("isAdmin") == "on" {
		isAdmin = 1
	}
	numId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	err = db.UpdateUserById(numId, name, surname, email, password, isAdmin) // ПРОДОЛЖИТЬ АПДЕЙТ ЮЗЕРОВ
	if err != nil {
		log.Println("Error updatiing user : ", err.Error())
		http.Error(w, "Internal server error", 500)
		return
	}
	log.Println("Edited user successfully")
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

type Message struct {
	Message string `json:"message"`
}

func MessageHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Читаем тело запроса
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusInternalServerError)
		return
	}

	// Декодируем JSON
	var message Message
	if err := json.Unmarshal(body, &message); err != nil {
		http.Error(w, "Unable to decode JSON", http.StatusBadRequest)
		return
	}

	// Получаем значение сообщения
	receivedMessage := message.Message
	fmt.Println("Received message:", receivedMessage)

	// Формируем URL для отправки запроса к Wit.ai
	url := "https://api.wit.ai/message?v=20240220&q=" + receivedMessage

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	witAccessToken := "UVJRI7BCS4FOHCMPRMPUZBXUM5GJBHC6"

	// Устанавливаем заголовок с токеном доступа к Wit.ai API
	req.Header.Set("Authorization", "Bearer "+witAccessToken)

	// Отправляем запрос к Wit.ai
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request to Wit.ai", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Декодируем ответ от Wit.ai
	var witResponse WitResponse
	if err := json.NewDecoder(resp.Body).Decode(&witResponse); err != nil {
		http.Error(w, "Failed to decode response from Wit.ai", http.StatusInternalServerError)
		return
	}
	log.Println("Response : ", witResponse)
	// Выводим текст ответа от Wit.ai
	fmt.Fprintf(w, "Wit.ai response: %s\n", witResponse.Text)

	// Отправляем ответ клиенту
	// w.WriteHeader(http.StatusOK) // Убираем эту строку, так как WriteHeader уже вызван в Fprintf
	// w.Write([]byte("Message received successfully")) // Эту строку тоже убираем
}

type WitResponse struct {
	MsgId string `json:"msg_id"`
	Text  string `json:"_text"`
}
