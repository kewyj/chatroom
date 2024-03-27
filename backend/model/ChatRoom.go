package model

import (
	"github.com/google/uuid"
)

// chatroom
type ChatRoom struct {
	ID       string            `json:"chatroom_id"`
	Users    map[string]string `json:"chat_users"`
	Messages []Message         `json:"messages"`
}

func NewChatRoom() ChatRoom {
	return ChatRoom{
		ID:       uuid.New().String(),
		Users:    make(map[string]string),
		Messages: []Message{},
	}
}
