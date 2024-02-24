package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/kewyj/chatroom/controller"
	"github.com/kewyj/chatroom/handler"
	"github.com/rs/cors"
)

func main() {
	cacheData := controller.NewChatService()
	handler := handler.NewHandler(cacheData)

	// routes
	r := mux.NewRouter()
	r.HandleFunc("/newuser", handler.NewUser).Methods("PUT")
	r.HandleFunc("/chat", handler.SendMessage).Methods("POST")
	r.HandleFunc("/poll", handler.Poll).Methods("PATCH")
	r.HandleFunc("/exit", handler.Exit).Methods("DELETE")

	// wrap with cors
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"PATCH", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"*"},
	})
	corsHandler := c.Handler(r)

	err := http.ListenAndServe(":3333", corsHandler)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
