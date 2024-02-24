package storage

import "github.com/kewyj/chatroom/model"

type Storage interface {
	AddNewChatRoom(cr model.ChatRoom) error
	AddNewUser(user model.User) error

	AddMessageToUser(userid string, msg model.Message) error

	GetUser(userid string) (model.User, error)
	GetUsers() ([]model.User, error)

	GetUserQueue(userid string) (model.MessageQueue, error)

	GetRoom(roomid string) (model.ChatRoom, error)
	GetRooms() ([]model.ChatRoom, error)

	RemoveUser(userid string) error
	RemoveRoom(roomid string) error
	ClearUserQueue(userid string) error
}
