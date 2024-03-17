package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kewyj/chatroom/controller"
	"github.com/kewyj/chatroom/handler"
	// "errors"
	// "fmt"
	// "net/http"
	// "os"
	//"github.com/gorilla/mux"
	//"github.com/rs/cors"
	//"github.com/aws/aws-lambda-go/events"
)

func main() {
	controller := controller.NewChatService()
	handler := handler.NewLambdaWrapper(controller)

	lambda.Start(handler.LambdaHandler)
	// cacheData := controller.NewChatService()
	// handler := handler.NewHandler(cacheData)

	// // routes
	// r := mux.NewRouter()
	// r.HandleFunc("/newuser", handler.NewUser).Methods("PUT")
	// r.HandleFunc("/chat", handler.SendMessage).Methods("POST")
	// r.HandleFunc("/poll", handler.Poll).Methods("PATCH")
	// r.HandleFunc("/exit", handler.Exit).Methods("DELETE")

	// // wrap with cors
	// c := cors.New(cors.Options{
	// 	AllowedOrigins: []string{"*"},
	// 	AllowedMethods: []string{"PATCH", "POST", "PUT", "DELETE"},
	// 	AllowedHeaders: []string{"*"},
	// })
	// corsHandler := c.Handler(r)

	// err := http.ListenAndServe(":3333", corsHandler)
	// if errors.Is(err, http.ErrServerClosed) {
	// 	fmt.Printf("server closed\n")
	// } else if err != nil {
	// 	fmt.Printf("error starting server: %s\n", err)
	// 	os.Exit(1)
	// }
}
