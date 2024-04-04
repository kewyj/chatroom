package controller

import "github.com/kewyj/chatroom/model"

type Controller interface {
	GetRooms() ([]model.GetRoomsResponse, error)
	AddRoom() (string, error)
	AddUser(user model.NewUserRequest) (string, error)
	AddUserToRoom(model.AddRoomRequest) error
	SendMessage(msg model.MessageRequest) error
	Poll(model.PollRequest) ([]model.Message, error)
	RemoveUserFromRoom(model.ExitRoomRequest) error
	RemoveUser(model.ExitRequest) error
	ClearAll(string) error
	Quit(uuid string, chatroom_id string) error
}
