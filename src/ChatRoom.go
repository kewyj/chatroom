package src

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

const MAX_USERS_IN_ROOM = 10
const MAX_ROOMS = 10
const INITIAL_ROOMS = 1

// chatroom
type ChatRoom struct {
	m_id       string
	m_users    []string
	m_messages *MessageQueue
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		m_id:       uuid.New().String(),
		m_users:    []string{},
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
	cr.SendSystemMessage(user[:4] + " has joined the room.")
}

func (cr *ChatRoom) GetMessages() *MessageQueue {
	return cr.m_messages
}

func (cr *ChatRoom) SendMessage(msg Message) {
	if cr.m_messages.Size() > MAX_MESSAGES_IN_ROOM {
		cr.m_messages.Dequeue()
	}
	cr.m_messages.Enqueue(msg)
}

func (cr *ChatRoom) SendSystemMessage(msg string) {
	if cr.m_messages.Size() > MAX_MESSAGES_IN_ROOM {
		cr.m_messages.Dequeue()
	}
	cr.m_messages.Enqueue(Message{".:system:.", msg})
}

type ChatService struct {
	// map of room name to Chatroom object
	m_rooms map[string]*ChatRoom
	// map of username to Chatroom object
	m_users map[string]*ChatRoom
	// map of subscribers to unread messages
	m_subs map[string]*MessageQueue
	// mutex
	mu sync.Mutex
}

func NewChatService() *ChatService {
	rooms := make(map[string]*ChatRoom)
	for i := 0; i < INITIAL_ROOMS; i++ {
		room := NewChatRoom()
		rooms[room.GetID()] = room
	}
	users := make(map[string]*ChatRoom)
	subs := make(map[string]*MessageQueue)
	return &ChatService{
		m_rooms: rooms,
		m_users: users,
		m_subs:  subs,
	}
}

func (cs *ChatService) PrintServerStatus() {
	fmt.Println("ROOMS AND USERS")
	for key, val := range cs.m_rooms {
		fmt.Println("ROOM: " + key[:4])
		fmt.Println("")
		fmt.Println("USERS")
		for _, user := range val.GetUsers() {
			fmt.Println(user[:4])
		}
		fmt.Println("")
		fmt.Println("MESSAGES")
		for _, msg := range *val.GetMessages() {
			fmt.Println(msg.Username, msg.Content)
		}
		fmt.Println("")
	}
	fmt.Println("")
	fmt.Println("QUEUED MESSAGES")
	for key, val := range cs.m_subs {
		fmt.Println("USER: " + key[:4])
		for _, msg := range *val {
			fmt.Println("	FROM: " + msg.Username[:4])
			fmt.Println("	MESSAGE: " + msg.Content)
		}
		fmt.Println("")
	}
}

func (cs *ChatService) AddUser() (string, error) {
	// count num rooms
	cs.mu.Lock()
	defer cs.mu.Unlock()
	count := 0
	name := uuid.New().String()
	for _, val := range cs.m_rooms {
		count++
		if val.GetCountUsers() < MAX_USERS_IN_ROOM {
			// found available room
			val.AddUser(name)
			cs.m_users[name] = val
			cs.m_subs[name] = &MessageQueue{}
			cs.updateSubscribers(val, Message{".:system:.", name[:4] + " has joined the room."})
			cs.PrintServerStatus()
			return name, nil
		}
	}
	if count >= MAX_ROOMS {
		// did not find available room + num rooms > MAX_ROOMS
		return "", errors.New("ChatService::AddUser - Rooms are at max capacity")
	}
	// did not find available room, create new chat room
	room := NewChatRoom()
	cs.m_rooms[room.GetID()] = room
	room.AddUser(name)
	cs.m_users[name] = room
	cs.m_subs[name] = &MessageQueue{}
	cs.updateSubscribers(room, Message{".:system:.", name[:4] + " has joined the room."})
	cs.PrintServerStatus()
	return name, nil
}

func (cs *ChatService) SendMessage(msg Message) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	val, ok := cs.m_users[msg.Username]
	if ok {
		// user exists
		msg.Username = msg.Username[:4] // shorten display name
		val.SendMessage(msg)
		cs.updateSubscribers(val, msg)
		cs.PrintServerStatus()
		return nil
	}
	return errors.New("ChatService::SendMessage - User not found")
}

func (cs *ChatService) SendSystemMessage(room string, msg string) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	val, ok := cs.m_rooms[room]
	if ok {
		val.SendSystemMessage(msg)
		cs.updateSubscribers(val, Message{".:system:.", msg})
		return nil
	}
	return errors.New("ChatService::SendSystemMessage - Room not found")
}

func (cs *ChatService) SendSystemMessageGlobal(msg string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	for _, val := range cs.m_rooms {
		val.SendSystemMessage(msg)
	}
}

func (cs *ChatService) GetRoom(room string) *ChatRoom {
	return cs.m_rooms[room]
}

func (cs *ChatService) updateSubscribers(cr *ChatRoom, msg Message) {
	for _, user := range cr.m_users {
		cs.m_subs[user].Enqueue(msg)
	}
}

func (cs *ChatService) RetrieveUndelivered(user string) *MessageQueue {
	return cs.m_subs[user]
}

func (cs *ChatService) RemoveUser(user string) error {
	// remove from chatroom
	room, ok := cs.m_users[user]
	if !ok {
		return errors.New("user not found")
	}
	index := -1
	for i, val := range room.m_users {
		if val == user {
			index = i
			break
		}
	}
	if index != -1 {
		room.m_users = append(room.m_users[:index], room.m_users[index+1:]...)
	}
	delete(cs.m_users, user)
	delete(cs.m_subs, user)
	// if room is empty delete it
	if room.IsEmpty() {
		delete(cs.m_rooms, room.m_id)
	} else {
		cs.SendSystemMessage(room.m_id, user[:4]+"has left the room.")
	}
	return nil
}
