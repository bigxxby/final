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

	"github.com/gorilla/websocket"
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
	boolean, user := isAuthenticated(r)
	if !boolean {
		user = &myDatabase.User{IsLogged: false}
	} else {
		user.IsLogged = true
	}
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
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	fmt.Println("FULL: ", string(body))
	fmt.Println("Intents: ", response.Intents)
	fmt.Println("Text:", response.Text)
	fmt.Println("Greetings Confidence:", response.Traits.Greetings)
	fmt.Println("Sentiment Value:", response.Traits.Sentiment)
	fmt.Println("Bye Value:", response.Traits.Bye)
	value := CheckTraits(&response) // ЭТО СООБЗЕНИЕ НАДО ОТПРАВИТЬ НА ФРОНТ

	if value == "404" {
		value = "I'm sorry, I didn't understand your request. Please Contact our Support team."
	} else if value == "negative" {
		value = "I understand that you're feeling negative about this. Let's see if I can help you in any other way."
	} else if value == "greetings" {
		value = "Hello there! Welcome to our platform. How can I assist you today?"

	} else if value == "bye" {
		value = "Thank you for chatting with me. If you have any more questions or need assistance in the future, feel free to reach out"

	} else if value == "thanks" {
		value = "You're welcome! If you have any more questions or need further assistance, don't hesitate to ask. Have a wonderful day!"

	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// пример JSON-структуры, которую вы хотите отправить
	response2 := struct {
		Message string `json:"message"`
	}{
		Message: value,
	}

	// кодируем JSON и отправляем его обратно на клиент
	json.NewEncoder(w).Encode(response2)
}

func CheckTraits(response *Response) string {
	if len(response.Traits.Sentiment) > 0 {
		for _, sentiment := range response.Traits.Sentiment {
			if sentiment.Value == "negative" {

				return "negative"
			}
		}
	}
	if len(response.Traits.Greetings) > 0 {
		for _, sentiment := range response.Traits.Greetings {
			if sentiment.Value == "true" && sentiment.Confidence > 0.9 {
				return "greetings"
			}
		}
	}
	if len(response.Traits.Thanks) > 0 {
		for _, sentiment := range response.Traits.Thanks {
			if sentiment.Value == "true" && sentiment.Confidence > 0.9 {
				return "thanks"
			}
		}
	}
	if len(response.Traits.Bye) > 0 {
		for _, sentiment := range response.Traits.Bye {
			if sentiment.Value == "true" && sentiment.Confidence > 0.9 {
				return "bye"
			}
		}
	}
	return "404"

}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HandleFunc(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		// Прочитать сообщение от клиента
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}

		// Вывести полученное сообщение на серверной стороне
		fmt.Printf("Received message: %s\n", p)

		// Отправить сообщение обратно клиенту
		err = conn.WriteMessage(messageType, p)
		if err != nil {
			fmt.Println("Write error:", err)
			break
		}
	}
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Отправка сообщения на клиент при необходимости
	err = conn.WriteMessage(websocket.TextMessage, []byte("Привет, клиент!"))
	if err != nil {
		log.Println(err)
		return
	}
	// upgrader.Upgrade(w, )
}

type Response struct {
	Entities map[string]interface{} `json:"entities"`
	Intents  []interface{}          `json:"intents"`
	Text     string                 `json:"text"`
	Traits   struct {
		Greetings []Trait `json:"wit$greetings"`
		Bye       []Trait `json:"wit$bye"`
		Sentiment []Trait `json:"wit$sentiment"`
		Thanks    []Trait `json:"wit$thanks"`
	} `json:"traits"`
}
type Trait struct {
	Confidence float64 `json:"confidence"`
	ID         string  `json:"id"`
	Value      string  `json:"value"`
}
type Intents struct {
	Confidence float64 `json:"confidence"`
	ID         string  `json:"id"`
	Name       string  `json:"name"`
}

func ChordHandler(w http.ResponseWriter, r *http.Request) {
	data, err := template.ParseFiles("ui/static/templates/chordPanel.html")
	if err != nil {
		log.Println(err)
		return
	}
	data.Execute(w, nil)
}

func GuitarHandler(w http.ResponseWriter, r *http.Request) {
	data, err := template.ParseFiles("ui/static/templates/guitar.html")
	if err != nil {
		log.Println(err)
		return
	}
	data.Execute(w, nil)
}
func BeginnerHandler(w http.ResponseWriter, r *http.Request) {
	data, err := template.ParseFiles("ui/static/templates/beginner.html")
	if err != nil {
		log.Println(err)
		return
	}
	data.Execute(w, nil)
}
func IntermediateHandler(w http.ResponseWriter, r *http.Request) {
	data, err := template.ParseFiles("ui/static/templates/intermediate.html")
	if err != nil {
		log.Println(err)
		return
	}
	data.Execute(w, nil)
}
func AdvancedHandler(w http.ResponseWriter, r *http.Request) {
	data, err := template.ParseFiles("ui/static/templates/advanced.html")
	if err != nil {
		log.Println(err)
		return
	}
	data.Execute(w, nil)
}
func MetronomeHandler(w http.ResponseWriter, r *http.Request) {
	data, err := template.ParseFiles("ui/static/templates/metronome.html")
	if err != nil {
		log.Println(err)
		return
	}
	data.Execute(w, nil)
}

type WitResponse struct {
	MsgId string `json:"msg_id"`
	Text  string `json:"_text"`
}
