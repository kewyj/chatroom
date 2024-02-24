package controller

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/kewyj/chatroom/limiter"
	"github.com/kewyj/chatroom/model"
	"github.com/kewyj/chatroom/storage"
)

type ChatService struct {
	storage  storage.Storage
	limiters map[string]*limiter.RateLimiter
}

func NewChatService() *ChatService {
	cache := storage.NewCache()
	for i := 0; i < model.INITIAL_ROOMS; i++ {
		cache.AddNewChatRoom(model.NewChatRoom())
	}

	return &ChatService{
		storage:  cache,
		limiters: make(map[string]*limiter.RateLimiter),
	}
}

func (cs *ChatService) PrintServerStatus() {
	fmt.Println("QUEUED MESSAGES")

	users, err := cs.storage.GetUsers()
	if err != nil {
		fmt.Println("Error getting users: %v", err.Error())
	}

	for _, user := range users {
		fmt.Println("USER: " + user.Username[:4])

		queue, err := cs.storage.GetUserQueue(user.Username)
		if err != nil {
			fmt.Println("error getting user queue: %v", err)
			continue
		}

		for _, msg := range queue {
			fmt.Println("	FROM: " + msg.Username[:4])
			fmt.Println("	MESSAGE: " + msg.Content)
		}
		fmt.Println("")
	}
}

func (cs *ChatService) AddUser() (string, error) {
	rooms, err := cs.storage.GetRooms()
	if err != nil {
		return "", fmt.Errorf("error getting rooms: %w", err)
	}

	for _, val := range rooms {
		if len(val.Users) < model.MAX_USERS_IN_ROOM {
			// found available room
			name, err := cs.addNewUser(val.ID)
			if err != nil {
				return "", fmt.Errorf("error creating user: %w", err)
			}

			return name, nil
		}
	}

	if len(rooms) >= model.MAX_ROOMS {
		// did not find available room + num rooms > MAX_ROOMS
		return "", errors.New("ChatService::AddUser - Rooms are at max capacity")
	}

	// did not find available room, create new chat room
	room := model.NewChatRoom()
	if err := cs.storage.AddNewChatRoom(room); err != nil {
		return "", fmt.Errorf("error creating chatroom: %w", err)
	}

	name, err := cs.addNewUser(room.ID)
	if err != nil {
		return "", fmt.Errorf("error creating user: %w", err)
	}

	return name, nil
}

func (cs *ChatService) addNewUser(roomID string) (string, error) {
	newUser := model.User{
		Username:   uuid.New().String(),
		ChatroomID: roomID,
	}

	if err := cs.storage.AddNewUser(newUser); err != nil {
		return "", err
	}

	cs.limiters[newUser.Username] = limiter.NewRateLimiter(model.MAX_MESSAGES_PER_SECOND)

	cs.updateSubscribers(roomID, model.Message{".:system:.", newUser.Username[:4] + " has entered the room."})

	return newUser.Username, nil
}

func (cs *ChatService) SendMessage(msg model.Message) error {
	user, err := cs.storage.GetUser(msg.Username)
	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	// user exists
	msg.Username = msg.Username[:4] // shorten display name

	cs.updateSubscribers(user.ChatroomID, msg)

	return nil
}

func (cs *ChatService) updateSubscribers(roomid string, msg model.Message) {
	room, err := cs.storage.GetRoom(roomid)
	if err != nil {
		fmt.Println("error getting room for update subscriber: %v", err)
		return
	}

	for _, user := range room.Users {
		if err := cs.storage.AddMessageToUser(user, msg); err != nil {
			fmt.Println("error updating user for update subscriber - %v: %v", user, err)
		}
	}
	
	cs.PrintServerStatus()
}

func (cs *ChatService) Poll(userid string) (model.MessageQueue, error) {
	queue, err := cs.storage.GetUserQueue(userid)
	if err != nil {
		return model.MessageQueue{}, fmt.Errorf("error getting user queue: %w", err)
	}

	if err := cs.storage.ClearUserQueue(userid); err != nil {
		return model.MessageQueue{}, fmt.Errorf("error clearing user queue: %w", err)
	}

	return queue, nil
}

func (cs *ChatService) RemoveUser(userid string) error {
	user, err := cs.storage.GetUser(userid)
	if err != nil {
		return fmt.Errorf("error getting user: %w", err)
	}

	if err := cs.storage.RemoveUser(userid); err != nil {
		return fmt.Errorf("error removing user: %w", err)
	}

	_, ok := cs.limiters[userid]
	if ok {
		cs.limiters[userid].Destroy()
		delete(cs.limiters, userid)
	}

	room, err := cs.storage.GetRoom(user.ChatroomID)
	if err != nil {
		return fmt.Errorf("error getting room: %w", err)
	}

	if len(room.Users) == 0 {
		err := cs.storage.RemoveRoom(room.ID)
		if err != nil {
			return fmt.Errorf("error removing room: %w", err)
		}

		return nil
	}

	cs.updateSubscribers(room.ID, model.Message{".:system:.", userid[:4] + " has left the room."})

	return nil
}

func (cs *ChatService) IsUserSpamming(userid string) bool {
	_, ok := cs.limiters[userid]
	if !ok {
		fmt.Println("user not found in is spamming")
		return true
	}

	select {
	case <-cs.limiters[userid].TokenBucket:
		return false
	default:
		return true
	}
}
