package model

import (
	"github.com/google/uuid"
)

// chatroom
type ChatRoom struct {
	ID       string
	Users    map[string]string
	Messages []Message
}

func NewChatRoom() ChatRoom {
	return ChatRoom{
		ID:       uuid.New().String(),
		Users:    make(map[string]string),
		Messages: []Message{},
	}
}
