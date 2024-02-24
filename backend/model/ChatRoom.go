package model

import (
	"github.com/google/uuid"
)

const MAX_USERS_IN_ROOM = 10
const MAX_ROOMS = 10
const INITIAL_ROOMS = 1
const MAX_MESSAGES_PER_SECOND = 2

// chatroom
type ChatRoom struct {
	ID    string
	Users []string
}

func NewChatRoom() ChatRoom {
	return ChatRoom{
		ID:    uuid.New().String(),
		Users: []string{},
	}
}
