package storage

import "github.com/kewyj/chatroom/model"

type Storage interface {
	NewChatRoom(cr model.ChatRoom) error
	AddUserToChatRoom(custom_username string, uuid string, chatroom_id string) error
	AddMessageToChatRoom(chatroom_id string, msg model.Message) error

	CheckIfRoomExists(chatroom_id string) bool
	GetRooms() ([]model.ChatRoom, error)
	GetUsername(chatroom_id string, uuid string) (string, error)
	GetRoomUsernames(chatroom_id string) ([]string, error)
	GetRoomUserUUIDs(chatroom_id string) ([]string, error)
	GetRoomMessages(chatroom_id string) ([]model.Message, error)

	RemoveEarliestMessage(chatroom_id string) error
	RemoveUserFromChatRoom(uuid string, chatroom_id string) error
	RemoveRoom(chatroom_id string) error
}
