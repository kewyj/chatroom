package model

import (
	"strings"

	"github.com/nwtgck/go-fakelish"
)

// chatroom
type ChatRoom struct {
	ID        string    `json:"chatroom_id"`
	UserCount int       `json:"user_count"`
	Messages  []Message `json:"messages"`
}

func NewChatRoom() ChatRoom {
	name := fakelish.GenerateFakeWord(6, 9)
	name = strings.Title(name)
	return ChatRoom{
		ID:       name,
		Messages: []Message{},
	}
}
