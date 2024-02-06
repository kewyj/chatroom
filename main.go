package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"encoding/json"
	"io/ioutil"
)

// Define a struct that matches the structure of the JSON data
type Message struct {
	Username string `json:"username"`
	Content    string `json:"content"`
}

type Queue []Message

func (q *Queue) Enqueue(item Message) {
    *q = append(*q, item)
}

func (q *Queue) Dequeue() Message {
    if len(*q) == 0 {
        return Message{}
    }
    item := (*q)[0]
    *q = (*q)[1:]
    return item
}

func (q *Queue) Size() int {
	return len(*q)
}

// map of rooms to list of users in that room
var rooms = make(map[string][]string)

// which user in which room
var users = make(map[string]string)

// map of rooms to list of messages in that room
var roomChat = make(map[string]*Queue)

func printRequest(r *http.Request) {
	fmt.Println("==============HTTP Request==============")
	// Accessing HTTP method
    method := r.Method
    fmt.Println("Method:", method)

    // Accessing URL
    url := r.URL
    fmt.Println("URL:", url)

    // Accessing HTTP headers
    headers := r.Header
    fmt.Println("Headers:", headers)

    // Accessing request body
    // Note: You typically read the request body when needed.
    // This is just an example to show how to access it.
    body := r.Body
    fmt.Println("Body:", body)

    // Accessing host
    host := r.Host
    fmt.Println("Host:", host)

    // Accessing remote address
    remoteAddr := r.RemoteAddr
    fmt.Println("RemoteAddr:", remoteAddr)
}

func printServerStatus() {
	fmt.Println("Rooms and Users")
	fmt.Println(rooms)
	fmt.Println("Users and Rooms")
	fmt.Println(users)
	fmt.Println("Rooms and Messages")
	fmt.Println(roomChat)
}

func addSystemMessage(room string, message string) {
	if (roomChat[room].Size() > 10) {
		roomChat[room].Dequeue()
	}
	roomChat[room].Enqueue(Message{"system", message})
}

func addMessage(msg Message) {
	room := users[msg.Username]
	if (roomChat[room].Size() > 10) {
		roomChat[room].Dequeue()
	}
	roomChat[room].Enqueue(msg)
}

func newUser(w http.ResponseWriter, r *http.Request) {

	// create name for user
	name := "user" + uuid.New().String()
	foundRoom := false
	roomCounter := 0

	for key, value := range	rooms {
		if len(value) < 10 {
			rooms[key] = append(rooms[key], name)
			users[name] = key
			foundRoom = true
			break
		}
		roomCounter++
		// limit num rooms to 10
		if (roomCounter > 10) {
			http.Error(w, "Rooms are full!", http.StatusServiceUnavailable)
			return
		}
	}

	if foundRoom == false {
		newRoom := uuid.New().String()
		rooms[newRoom] = []string{ name }
		users[name] = newRoom
	}

	// set session token
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   name,
		MaxAge: -1,
	})

	io.WriteString(w, "Welcome " + name + "!")
	printRequest(r)

	w.WriteHeader(http.StatusOK) // redirect to chatroom
	printServerStatus()
}

func chat(w http.ResponseWriter, r *http.Request) {

	


	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var msg Message
	
	if err = json.Unmarshal(body, &msg); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	room, ok := users[msg.Username]
	if !ok {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}
	if _, ok := roomChat[room]; !ok {
		roomChat[room] = &Queue{}
	}
	roomChat[room].Enqueue(msg)

	fmt.Printf("new message: %s", msg.Content)
	printRequest(r)
	w.WriteHeader(http.StatusOK)
	printServerStatus()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/newuser", newUser).Methods("PUT")
	r.HandleFunc("/chat", chat).Methods("POST")

	err := http.ListenAndServe(":3333", r)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}