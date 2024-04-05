package storage

import "github.com/kewyj/chatroom/model"

type Storage interface {
	Initialize() error

	CheckIfRoomExists(chatroom_id string) bool

	NewUser(uuid string, username string) error
	NewChatRoom(id string) error

	AddUserToChatRoom(chatroom_id string) error
	AddMessageToChatRoom(chatroom_id string, msg model.Message) error

	GetRooms() ([]model.ChatRoom, error)
	GetRoom(chatroom_id string) (model.ChatRoom, error)
	GetUsername(uuid string) (string, error)
	GetRoomMessages(chatroom_id string) ([]model.Message, error)

	RemoveEarliestMessage(chatroom_id string) error
	RemoveUserFromChatRoom(chatroom_id string) error

	RemoveUser(uuid string) error
	RemoveRoom(chatroom_id string) error

	UpdateUserActivity(uuid string, time string) error

	CullIfLessThan(time string) error
	ClearAll() error
}
