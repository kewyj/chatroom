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

	case "/rooms":
		if request.RequestContext.HTTP.Method == "GET" {
			return l.Rooms(request)
		}

	case "/newroom":
		if request.RequestContext.HTTP.Method == "PUT" {
			return l.NewRoom(request)
		}

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

func (l *LambdaHandler) Rooms(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Println("Handling Rooms Request")

	rooms, err := l.controller.GetRooms()
	if err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: 500}, err
	}

	roomsJSON, err := json.Marshal(rooms)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: 500}, err
	}

	response := events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       string(roomsJSON),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return response, nil
}

func (l *LambdaHandler) NewRoom(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Println("Handling New Room Request")

	roomid, err := l.controller.AddRoom()
	if err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: 500}, err
	}

	roomidResponse := model.NewRoomResponse{
		RoomID: roomid,
	}

	roomidJSON, err := json.Marshal(roomidResponse)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: 500}, err
	}

	response := events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       string(roomidJSON),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return response, nil
}

func (l *LambdaHandler) NewUser(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Println("Handling New User Request")

	user, err := l.UnmarshalNewUser(request.Body)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: 500}, err
	}

	uuid, err := l.controller.AddUser(user)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: 500}, err
	}

	uuidResponse := model.NewUserResponse{
		Username: uuid,
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

	msg, err := l.UnmarshalMessageRequest(request.Body)
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

	pollReq, err := l.UnmarshalPollRequest(request.Body)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: 500}, err
	}

	messages, err := l.controller.Poll(pollReq)
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

	msg, err := l.UnmarshalExitRequest(request.Body)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: 500}, err
	}

	err = l.controller.RemoveUser(msg)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: 500}, err
	}

	return events.APIGatewayV2HTTPResponse{StatusCode: 200}, nil
}

func (l *LambdaHandler) UnmarshalNewUser(body string) (model.NewUserRequest, error) {
	var user model.NewUserRequest
	err := json.Unmarshal([]byte(body), &user)
	if err != nil {
		return model.NewUserRequest{}, err
	}

	return user, nil
}

func (l *LambdaHandler) UnmarshalMessageRequest(body string) (model.MessageRequest, error) {
	var msg model.MessageRequest
	err := json.Unmarshal([]byte(body), &msg)
	if err != nil {
		return model.MessageRequest{}, err
	}

	return msg, nil
}

func (l *LambdaHandler) UnmarshalPollRequest(body string) (model.PollRequest, error) {
	var msg model.PollRequest
	err := json.Unmarshal([]byte(body), &msg)
	if err != nil {
		return model.PollRequest{}, err
	}

	return msg, nil
}

func (l *LambdaHandler) UnmarshalExitRequest(body string) (model.ExitRequest, error) {
	var msg model.ExitRequest
	err := json.Unmarshal([]byte(body), &msg)
	if err != nil {
		return model.ExitRequest{}, err
	}

	return msg, nil
}
