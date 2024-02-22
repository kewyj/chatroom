package src

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type ChatService struct {
	// map of room name to Chatroom object
	m_rooms map[string]*ChatRoom

	// map of username to Chatroom object
	m_users map[string]*ChatRoom

	// map of subscribers to unread messages
	m_subs map[string]*MessageQueue

	// map of username to RateLimiter
	m_limiter map[string]*RateLimiter

	// mutex
	mu sync.Mutex
}

func NewChatService() *ChatService {
	rooms := make(map[string]*ChatRoom)
	for i := 0; i < INITIAL_ROOMS; i++ {
		room := NewChatRoom()
		rooms[room.GetID()] = room
	}

	return &ChatService{
		m_rooms:   rooms,
		m_users:   make(map[string]*ChatRoom),
		m_subs:    make(map[string]*MessageQueue),
		m_limiter: make(map[string]*RateLimiter),
	}
}

func (cs *ChatService) PrintServerStatus() {
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
			cs.addNewUser(name, val)
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
	cs.addNewUser(name, room)
	return name, nil
}

func (cs *ChatService) addNewUser(name string, room *ChatRoom) {
	cs.m_users[name] = room
	cs.m_subs[name] = &MessageQueue{}
	cs.m_limiter[name] = NewRateLimiter(MAX_MESSAGES_PER_SECOND)

	room.AddUser(name)
	cs.updateSubscribers(room, Message{".:system:.", name[:4] + " has joined the room."})
	cs.PrintServerStatus()
}

func (cs *ChatService) SendMessage(msg Message) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	val, ok := cs.m_users[msg.Username]

	if !ok {
		return errors.New("ChatService::SendMessage - User not found")
	}

	// user exists
	msg.Username = msg.Username[:4] // shorten display name
	cs.updateSubscribers(val, msg)
	cs.PrintServerStatus()

	return nil
}

func (cs *ChatService) SendSystemMessage(room string, msg string) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	val, ok := cs.m_rooms[room]

	if !ok {
		return errors.New("ChatService::SendSystemMessage - Room not found")
	}

	cs.updateSubscribers(val, Message{".:system:.", msg})
	return nil
}

func (cs *ChatService) SendSystemMessageGlobal(msg string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	for _, val := range cs.m_rooms {
		cs.updateSubscribers(val, Message{".:system:.", msg})
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

func (cs *ChatService) RetrieveUndelivered(user string) (*MessageQueue, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	val, ok := cs.m_subs[user]
	if !ok {
		return nil, errors.New("user not found")
	}
	return val, nil
}

func (cs *ChatService) RemoveUser(user string) error {
	// remove from chatroom
	room, ok := cs.m_users[user]
	if !ok {
		return errors.New("user not found")
	}

	cs.deleteUserFromRoom(user, room)
	cs.deleteUser(user)

	// if room is empty delete it
	if room.IsEmpty() {
		delete(cs.m_rooms, room.m_id)
	} else {
		cs.SendSystemMessage(room.m_id, user[:4]+" has left the room.")
		cs.PrintServerStatus()
	}
	return nil
}

func (cs *ChatService) deleteUserFromRoom(user string, room *ChatRoom) {
	// find room with user
	index := -1
	for i, val := range room.m_users {
		if val == user {
			index = i
			break
		}
	}

	// remove user from room
	if index == -1 {
		return
	}
	room.m_users = append(room.m_users[:index], room.m_users[index+1:]...)
}

func (cs *ChatService) deleteUser(user string) {
	// remove references to user
	delete(cs.m_users, user)
	delete(cs.m_subs, user)
	cs.m_limiter[user].Destroy()
	delete(cs.m_limiter, user)
}

func (cs *ChatService) IsUserSpamming(user string) bool {
	select {
	case <-cs.m_limiter[user].TokenBucket:
		return false
	default:
		return true
	}
}
