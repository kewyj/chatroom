package handler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/kewyj/chatroom/controller"
	"github.com/kewyj/chatroom/model"
)

type LambdaHandler struct {
	controller controller.Controller
}

func NewLambdaWrapper(c controller.Controller) *LambdaHandler {
	return &LambdaHandler{
		controller: c,
	}
}

func (l *LambdaHandler) LambdaHandler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Println("Request Path: ", request.RequestContext.HTTP.Path)
	log.Println("Request Method: ", request.RequestContext.HTTP.Method)

	switch request.RequestContext.HTTP.Path {

	case "/newuser":
		if request.RequestContext.HTTP.Method == "PUT" {
			return l.NewUser(request)
		}

	case "/chat":
		if request.RequestContext.HTTP.Method == "POST" {
			return l.Chat(request)
		}

	case "/poll":
		if request.RequestContext.HTTP.Method == "PATCH" {
			return l.Poll(request)
		}

	case "/exit":
		if request.RequestContext.HTTP.Method == "DELETE" {
			return l.Exit(request)
		}

	default:
		return events.APIGatewayV2HTTPResponse{StatusCode: 404}, nil
	}

	return events.APIGatewayV2HTTPResponse{StatusCode: 405}, nil
}

func (l *LambdaHandler) NewUser(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Println("Handling New User Request")

	uuid, err := l.controller.AddUser()
	if err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: 500}, err
	}

	uuidResponse := NewUserResponse{
		Username: &uuid,
	}

	uuidJSON, err := json.Marshal(uuidResponse)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: 500}, err
	}

	response := events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       string(uuidJSON),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return response, nil
}

func (l *LambdaHandler) Chat(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Println("Handling Chat Request")

	msg, err := l.UnmarshalMessage(request.Body)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: 500}, err
	}

	err = l.controller.SendMessage(msg)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: 500}, err
	}

	return events.APIGatewayV2HTTPResponse{StatusCode: 200}, nil
}

func (l *LambdaHandler) Poll(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Println("Handling Poll Request")

	msg, err := l.UnmarshalMessage(request.Body)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: 500}, err
	}

	messages, err := l.controller.Poll(msg.Username)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: 500}, err
	}

	messagesJSON, err := json.Marshal(messages)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: 500}, err
	}

	response := events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       string(messagesJSON),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return response, nil
}

func (l *LambdaHandler) Exit(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Println("Handling Exit Request")

	msg, err := l.UnmarshalMessage(request.Body)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: 500}, err
	}

	err = l.controller.RemoveUser(msg.Username)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: 500}, err
	}

	return events.APIGatewayV2HTTPResponse{StatusCode: 200}, nil
}

func (l *LambdaHandler) UnmarshalMessage(body string) (model.Message, error) {
	var msg model.Message
	err := json.Unmarshal([]byte(body), &msg)
	if err != nil {
		return model.Message{}, err
	}

	return msg, nil
}

type NewUserResponse struct {
	Username *string `json:"username"`
}
