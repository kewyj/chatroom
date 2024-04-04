package model

import (
	"github.com/google/uuid"
)

// chatroom
type ChatRoom struct {
	ID        string    `json:"chatroom_id"`
	UserCount int       `json:"user_count"`
	Messages  []Message `json:"messages"`
}

func NewChatRoom() ChatRoom {
	return ChatRoom{
		ID:       uuid.New().String(),
		Messages: []Message{},
	}
}
