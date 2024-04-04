package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

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

	if request.RequestContext.HTTP.Method == "OPTIONS" {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"Access-Control-Allow-Headers":     "Content-Type, X-AMZ-DATE, Authorization, X-Api-Key, X-Amz-Security-Token",
				"Access-Control-Allow-Origin":      "*",
				"Access-Control-Allow-Methods":     "OPTIONS, GET, POST, PUT, DELETE",
				"Access-Control-Allow-Credentials": "true",
			},
			Body: "",
		}, nil
	}

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

	case "/addtoroom":
		if request.RequestContext.HTTP.Method == "PUT" {
			return l.AddToRoom(request)
		}

	case "/chat":
		if request.RequestContext.HTTP.Method == "POST" {
			return l.Chat(request)
		}

	case "/poll":
		if request.RequestContext.HTTP.Method == "PATCH" {
			return l.Poll(request)
		}
	case "/exitroom":
		if request.RequestContext.HTTP.Method == "DELETE" {
			return l.ExitRoom(request)
		}

	case "/exit":
		if request.RequestContext.HTTP.Method == "DELETE" {
			return l.Exit(request)
		}

	case "/clear":
		if request.RequestContext.HTTP.Method == "DELETE" {
			return l.Clear(request)
		}

	case "/quit":
		if request.RequestContext.HTTP.Method == "DELETE" {
			return l.Quit(request)
		}

	default:
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 404,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 405,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func (l *LambdaHandler) Rooms(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Println("Handling Rooms Request")

	rooms, err := l.controller.GetRooms()
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		}, err
	}

	roomsJSON, err := json.Marshal(rooms)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		}, err
	}

	response := events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       string(roomsJSON),
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "*",
			"Access-Control-Allow-Headers": "*",
		},
	}

	return response, nil
}

func (l *LambdaHandler) NewRoom(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Println("Handling New Room Request")

	roomid, err := l.controller.AddRoom()
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		}, err
	}

	roomidResponse := model.NewRoomResponse{
		RoomID: roomid,
	}

	roomidJSON, err := json.Marshal(roomidResponse)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		}, err
	}

	response := events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       string(roomidJSON),
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "*",
			"Access-Control-Allow-Headers": "*",
		},
	}

	return response, nil
}

func (l *LambdaHandler) NewUser(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Println("Handling New User Request")

	user, err := l.UnmarshalNewUser(request.Body)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		}, err
	}

	uuid, err := l.controller.AddUser(user)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		}, err
	}

	uuidResponse := model.NewUserResponse{
		Username: uuid,
	}

	uuidJSON, err := json.Marshal(uuidResponse)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		}, err
	}

	response := events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       string(uuidJSON),
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "*",
			"Access-Control-Allow-Headers": "*",
		},
	}

	return response, nil
}

func (l *LambdaHandler) AddToRoom(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Println("Handling Add To Room Request")

	msg, err := l.UnmarshalAddRoomRequest(request.Body)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		}, err
	}

	err = l.controller.AddUserToRoom(msg)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		}, err
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "*",
			"Access-Control-Allow-Headers": "*",
		},
	}, nil
}

func (l *LambdaHandler) Chat(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Println("Handling Chat Request")

	msg, err := l.UnmarshalMessageRequest(request.Body)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		}, err
	}

	err = l.controller.SendMessage(msg)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		}, err
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "*",
			"Access-Control-Allow-Headers": "*",
		},
	}, nil
}

func (l *LambdaHandler) Poll(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Println("Handling Poll Request")

	pollReq, err := l.UnmarshalPollRequest(request.Body)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		}, err
	}

	messages, err := l.controller.Poll(pollReq)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		}, err
	}

	messagesJSON, err := json.Marshal(messages)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		}, err
	}

	response := events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       string(messagesJSON),
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "*",
			"Access-Control-Allow-Headers": "*",
		},
	}

	return response, nil
}

func (l *LambdaHandler) ExitRoom(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Println("Handling Exit Room Request")

	msg, err := l.UnmarshalExitRoomRequest(request.Body)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		}, err
	}

	err = l.controller.RemoveUserFromRoom(msg)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		}, err
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "*",
			"Access-Control-Allow-Headers": "*",
		},
	}, nil
}

func (l *LambdaHandler) Exit(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Println("Handling Exit Request")

	msg, err := l.UnmarshalExitRequest(request.Body)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		}, err
	}

	err = l.controller.RemoveUser(msg)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		}, err
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "*",
			"Access-Control-Allow-Headers": "*",
		},
	}, nil
}

func (l *LambdaHandler) Clear(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Println("Handling Clear Request")

	msg, err := l.UnmarshalClearRequest(request.Body)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		}, err
	}

	err = l.controller.ClearAll(msg.Password)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		}, err
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "*",
			"Access-Control-Allow-Headers": "*",
		},
	}, nil
}

func (l *LambdaHandler) Quit(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	msg, err := l.UnmarshalQuitRequest(request.Body)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		}, err
	}

	err = l.controller.Quit(msg.Username, msg.RoomID)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		}, err
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "*",
			"Access-Control-Allow-Headers": "*",
		},
	}, nil
}

func (l *LambdaHandler) UnmarshalNewUser(body string) (model.NewUserRequest, error) {
	var user model.NewUserRequest
	err := json.Unmarshal([]byte(body), &user)
	if err != nil {
		return model.NewUserRequest{}, err
	}

	return user, nil
}

func (l *LambdaHandler) UnmarshalAddRoomRequest(body string) (model.AddRoomRequest, error) {
	var msg model.AddRoomRequest
	err := json.Unmarshal([]byte(body), &msg)
	if err != nil {
		return model.AddRoomRequest{}, err
	}

	return msg, nil
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

func (l *LambdaHandler) UnmarshalExitRoomRequest(body string) (model.ExitRoomRequest, error) {
	var msg model.ExitRoomRequest
	err := json.Unmarshal([]byte(body), &msg)
	if err != nil {
		return model.ExitRoomRequest{}, err
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

func (l *LambdaHandler) UnmarshalClearRequest(body string) (model.ClearRequest, error) {
	var msg model.ClearRequest
	err := json.Unmarshal([]byte(body), &msg)
	if err != nil {
		return model.ClearRequest{}, err
	}

	return msg, nil
}

func (l *LambdaHandler) UnmarshalQuitRequest(body string) (model.QuitRequest, error) {
	var msg model.QuitRequest
	err := json.Unmarshal([]byte(body), &msg)
	if err != nil {
		return model.QuitRequest{}, err
	}

	return msg, nil
}
