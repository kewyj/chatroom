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
	// TestDynamo()

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

// func TestDynamo() {
// 	sess := session.Must(session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable}))

// 	svc := dynamodb.New(sess)

// 	input := &dynamodb.ListTablesInput{}

// 	fmt.Println("Tables:")

// 	for {
// 		// Get the list of tables
// 		result, err := svc.ListTables(input)
// 		if err != nil {
// 			if aerr, ok := err.(awserr.Error); ok {
// 				switch aerr.Code() {
// 				case dynamodb.ErrCodeInternalServerError:
// 					fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
// 				default:
// 					fmt.Println(aerr.Error())
// 				}
// 			} else {
// 				// Print the error, cast err to awserr.Error to get the Code and
// 				// Message from an error.
// 				fmt.Println(err.Error())
// 			}
// 			return
// 		}

// 		for _, n := range result.TableNames {
// 			fmt.Println(*n)
// 		}

// 		// assign the last read tablename as the start for our next call to the ListTables function
// 		// the maximum number of table names returned in a call is 100 (default), which requires us to make
// 		// multiple calls to the ListTables function to retrieve all table names
// 		input.ExclusiveStartTableName = result.LastEvaluatedTableName

// 		if result.LastEvaluatedTableName == nil {
// 			break
// 		}
// 	}
// }
