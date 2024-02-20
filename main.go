package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/kewyj/chatroom/src"
)

func main() {
	cacheData := src.NewChatService()
	handler := src.NewHandler(cacheData)

	r := mux.NewRouter()
	r.HandleFunc("/newuser", handler.NewUser).Methods("PUT")
	r.HandleFunc("/chat", handler.SendMessage).Methods("POST")
	r.HandleFunc("/poll", handler.Poll).Methods("GET")

	err := http.ListenAndServe(":3333", r)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
