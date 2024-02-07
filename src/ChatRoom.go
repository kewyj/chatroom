package src

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
)

const MAX_USERS_IN_ROOM = 10
const MAX_ROOMS = 10
const INITIAL_ROOMS = 1

// chatroom
type ChatRoom struct {
	m_id string
	m_users []string
	m_messages *MessageQueue
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		m_id: uuid.New().String(),
		m_users: []string{},
		m_messages: &MessageQueue{},
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
	cr.SendSystemMessage(user + " has joined the room.")
}

func (cr *ChatRoom) GetMessages() *MessageQueue {
	return cr.m_messages
}

func (cr *ChatRoom) SendMessage(msg Message) {
	cr.m_messages.Enqueue(msg)
}

func (cr *ChatRoom) SendSystemMessage(msg string) {
	cr.m_messages.Enqueue(Message{".:system:.", msg})
}

type ChatService struct {
	m_rooms map[string]*ChatRoom
	m_users map[string]*ChatRoom
}

func NewChatService() *ChatService {
	rooms := make(map[string]*ChatRoom)
	for i := 0; i < INITIAL_ROOMS; i++ {
		room := NewChatRoom()
		rooms[room.GetID()] = room
	}
	users := make(map[string]*ChatRoom)
	return &ChatService{
		m_rooms: rooms,
		m_users: users,
	}
}

func (cs *ChatService) PrintServerStatus() {
	fmt.Println("ROOMS AND USERS")
	for key, val := range cs.m_rooms {
		fmt.Println("ROOM: " + key)
		fmt.Println("")
		fmt.Println("USERS")
		for _, user := range val.GetUsers() {
			fmt.Println(user)
		}
		fmt.Println("")
		fmt.Println("MESSAGES")
		for _, msg := range *val.GetMessages() {
			fmt.Println(msg.Username)
			fmt.Println(msg.Content)
		}
		fmt.Println("")
	}
}

func (cs *ChatService) AddUser() (string, error) {
	count := 0
	name := uuid.New().String()
	for _, val := range cs.m_rooms {
		count++
		if val.GetCountUsers() < MAX_USERS_IN_ROOM {
			val.AddUser(name)
			cs.m_users[name] = val
			cs.PrintServerStatus()
			return name, nil
		}
	}
	if count >= MAX_ROOMS {
		return "", errors.New("ChatService::AddUser - Rooms are at max capacity!")
	}
	room := NewChatRoom()
	rooms[room.GetID()] = room
	room.AddUser(name)
	cs.m_users[name] = room
}

func (cs *ChatService) SendMessage(msg Message) error {
	val, ok := cs.m_users[msg.Username]
	if ok {
		val.SendMessage(msg)
		cs.PrintServerStatus()
		return nil
	}
	return errors.New("ChatService::SendMessage - User not found.")
}

func (cs *ChatService) SendSystemMessage(room string, msg string) error {
	val, ok := cs.m_rooms[room]
	if ok {
		val.SendSystemMessage(msg)
		return nil
	}
	return errors.New("ChatService::SendSystemMessage - Room not found.")
}

func (cs *ChatService) SendSystemMessageGlobal(msg string) {
	for _, val := range cs.m_rooms {
		val.SendSystemMessage(msg)
	}
}

func (cs *ChatService) GetRoom(room string) *ChatRoom {
	return cs.m_rooms[room]
}