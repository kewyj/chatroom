package src

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type User struct {
	chatroom *ChatRoom
	queue    *MessageQueue
	limiter  *RateLimiter
}

func NewUser() *User {
	return &User{}
}

type ChatService struct {
	// map of room name to Chatroom object
	rooms map[string]*ChatRoom

	// map of username to user object
	users map[string]*User

	// mutex
	mu sync.Mutex
}

func NewChatService() *ChatService {
	rms := make(map[string]*ChatRoom)
	for i := 0; i < INITIAL_ROOMS; i++ {
		room := NewChatRoom()
		rms[room.GetID()] = room
	}

	return &ChatService{
		rooms: rms,
		users: make(map[string]*User),
	}
}

func (cs *ChatService) PrintServerStatus() {
	fmt.Println("QUEUED MESSAGES")
	for key, val := range cs.users {
		fmt.Println("USER: " + key[:4])
		for _, msg := range *val.queue {
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

	name := uuid.New().String()

	for _, val := range cs.rooms {
		if val.GetCountUsers() < MAX_USERS_IN_ROOM {
			// found available room
			cs.addNewUser(name, val)
			return name, nil
		}
	}

	if len(cs.rooms) >= MAX_ROOMS {
		// did not find available room + num rooms > MAX_ROOMS
		return "", errors.New("ChatService::AddUser - Rooms are at max capacity")
	}

	// did not find available room, create new chat room
	room := NewChatRoom()
	cs.rooms[room.GetID()] = room
	cs.addNewUser(name, room)
	return name, nil
}

func (cs *ChatService) addNewUser(username string, room *ChatRoom) {
	cs.users[username] = NewUser()
	cs.users[username].chatroom = room
	cs.users[username].queue = &MessageQueue{}
	cs.users[username].limiter = NewRateLimiter(MAX_MESSAGES_PER_SECOND)

	room.AddUser(username)
	cs.updateSubscribers(room, Message{".:system:.", username[:4] + " has joined the room."})
	cs.PrintServerStatus()
}

func (cs *ChatService) SendMessage(msg Message) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	val, ok := cs.users[msg.Username]

	if !ok {
		return errors.New("ChatService::SendMessage - User not found")
	}

	// user exists
	msg.Username = msg.Username[:4] // shorten display name
	cs.updateSubscribers(val.chatroom, msg)
	cs.PrintServerStatus()

	return nil
}

func (cs *ChatService) SendSystemMessage(room string, msg string) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	val, ok := cs.rooms[room]

	if !ok {
		return errors.New("ChatService::SendSystemMessage - Room not found")
	}

	cs.updateSubscribers(val, Message{".:system:.", msg})
	return nil
}

func (cs *ChatService) SendSystemMessageGlobal(msg string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	for _, val := range cs.rooms {
		cs.updateSubscribers(val, Message{".:system:.", msg})
	}
}

func (cs *ChatService) GetRoom(room string) *ChatRoom {
	return cs.rooms[room]
}

func (cs *ChatService) updateSubscribers(cr *ChatRoom, msg Message) {
	for _, user := range cr.m_users {
		cs.users[user].queue.Enqueue(msg)
	}
}

func (cs *ChatService) RetrieveUndelivered(username string) (*MessageQueue, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	val, ok := cs.users[username]
	if !ok {
		return nil, errors.New("user not found")
	}
	return val.queue, nil
}

func (cs *ChatService) RemoveUser(username string) error {
	// remove from chatroom
	user, ok := cs.users[username]
	if !ok {
		return errors.New("user not found")
	}

	cs.deleteUserFromRoom(username, user.chatroom)
	cs.deleteUser(username)

	// if room is empty delete it
	if user.chatroom.IsEmpty() {
		delete(cs.rooms, user.chatroom.m_id)
	} else {
		cs.SendSystemMessage(user.chatroom.m_id, username[:4]+" has left the room.")
		cs.PrintServerStatus()
	}
	return nil
}

func (cs *ChatService) deleteUserFromRoom(username string, room *ChatRoom) {
	// find room with user
	index := -1
	for i, val := range room.m_users {
		if val == username {
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

func (cs *ChatService) deleteUser(username string) {
	// remove references to user
	cs.users[username].limiter.Destroy()
	delete(cs.users, username)
}

func (cs *ChatService) IsUserSpamming(username string) bool {
	select {
	case <-cs.users[username].limiter.TokenBucket:
		return false
	default:
		return true
	}
}
