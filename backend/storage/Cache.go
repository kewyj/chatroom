package storage

import (
	"errors"
	"sync"

	"github.com/kewyj/chatroom/model"
)

type Cache struct {
	// map of room name to Chatroom object
	rooms map[string]model.ChatRoom

	// mutex
	mu sync.Mutex
}

func NewCache() *Cache {
	return &Cache{
		rooms: make(map[string]model.ChatRoom),
	}
}

func (c *Cache) NewChatRoom(cr model.ChatRoom) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.rooms[cr.ID]
	if ok {
		return errors.New("tried to create chatroom that already exists")
	}

	c.rooms[cr.ID] = cr
	return nil
}

func (c *Cache) AddUserToChatRoom(custom_username string, uuid string, chatroom_id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	chatroom, ok := c.rooms[chatroom_id]
	if !ok {
		return errors.New("tried to add user to chatroom that does not exist")
	}

	chatroom.Users[uuid] = custom_username
	return nil
}

func (c *Cache) AddMessageToChatRoom(chatroom_id string, msg model.Message) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	chatroom, ok := c.rooms[chatroom_id]
	if !ok {
		return errors.New("tried to add message to chatroom that does not exist")
	}

	chatroom.Messages = append(chatroom.Messages, msg)
	return nil
}

func (c *Cache) CheckIfRoomExists(chatroom_id string) bool {
	_, ok := c.rooms[chatroom_id]
	return ok
}

func (c *Cache) GetRooms() ([]model.ChatRoom, error) {
	result := make([]model.ChatRoom, 0, len(c.rooms))

	for _, val := range c.rooms {
		result = append(result, val)
	}

	return result, nil
}

func (c *Cache) GetUsername(chatroom_id string, uuid string) (string, error) {
	chatroom, ok := c.rooms[chatroom_id]
	if !ok {
		return "", errors.New("tried to get usernames from chatroom that does not exist")
	}

	username, ok := chatroom.Users[uuid]
	if !ok {
		return "", errors.New("tried to get username with uuid that does not exist")
	}

	return username, nil
}

func (c *Cache) GetRoomUsernames(chatroom_id string) ([]string, error) {
	chatroom, ok := c.rooms[chatroom_id]
	if !ok {
		return []string{}, errors.New("tried to get usernames from chatroom that does not exist")
	}

	usernames := make([]string, 0, len(chatroom.Users))
	for _, username := range chatroom.Users {
		usernames = append(usernames, username)
	}

	return usernames, nil
}

func (c *Cache) GetRoomUserUUIDs(chatroom_id string) ([]string, error) {
	chatroom, ok := c.rooms[chatroom_id]
	if !ok {
		return []string{}, errors.New("tried to get user uuid from chatroom that does not exist")
	}

	user_uuids := make([]string, 0, len(chatroom.Users))
	for user_uuid := range chatroom.Users {
		user_uuids = append(user_uuids, user_uuid)
	}

	return user_uuids, nil
}

func (c *Cache) GetRoomMessages(chatroom_id string) ([]model.Message, error) {
	chatroom, ok := c.rooms[chatroom_id]
	if !ok {
		return []model.Message{}, errors.New("tried to get messages from chatroom that does not exist")
	}

	return chatroom.Messages, nil
}

func (c *Cache) RemoveEarliestMessage(chatroom_id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	chatroom, ok := c.rooms[chatroom_id]
	if !ok {
		return errors.New("tried to remove message from chatroom that does not exist")
	}

	if len(chatroom.Messages) == 0 {
		return nil
	}

	chatroom.Messages = chatroom.Messages[1:]
	return nil
}

func (c *Cache) RemoveUserFromChatRoom(uuid string, chatroom_id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	chatroom, ok := c.rooms[chatroom_id]
	if !ok {
		return errors.New("tried to remove user from chatroom that does not exist")
	}

	_, ok = chatroom.Users[uuid]
	if !ok {
		return errors.New("tried to remove user that does not exist")
	}

	delete(chatroom.Users, uuid)
	return nil
}

func (c *Cache) RemoveRoom(chatroom_id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.rooms[chatroom_id]
	if !ok {
		return errors.New("tried to remove chatroom that does not exist")
	}

	delete(c.rooms, chatroom_id)
	return nil
}
