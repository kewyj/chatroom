package src

import (
	"github.com/google/uuid"
)

const MAX_USERS_IN_ROOM = 10
const MAX_ROOMS = 10
const INITIAL_ROOMS = 1
const MAX_MESSAGES_PER_SECOND = 2

// chatroom
type ChatRoom struct {
	m_id    string
	m_users []string
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		m_id:    uuid.New().String(),
		m_users: []string{},
	}
}

func (cr *ChatRoom) GetID() string {
	return cr.m_id
}

func (cr *ChatRoom) GetUsers() []string {
	return cr.m_users
}

func (cr *ChatRoom) GetCountUsers() int {
	return len(cr.m_users)
}

func (cr *ChatRoom) IsEmpty() bool {
	return len(cr.m_users) == 0
}

func (cr *ChatRoom) AddUser(user string) {
	cr.m_users = append(cr.m_users, user)
}
