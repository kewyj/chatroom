package controller

import "github.com/kewyj/chatroom/model"

type Controller interface {
	GetRooms() ([]model.GetRoomsResponse, error)
	AddUser(user model.NewUserRequest) (string, error)
	SendMessage(msg model.MessageRequest) error
	Poll(model.PollRequest) ([]model.Message, error)
	RemoveUser(model.ExitRequest) error
}
