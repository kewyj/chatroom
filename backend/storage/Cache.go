package storage

import (
	"errors"
	"sync"

	"github.com/kewyj/chatroom/model"
)

type Cache struct {
	// map of room name to Chatroom object
	rooms map[string]model.ChatRoom

	// map of username to user object
	users map[string]model.User

	// map of username to queue object
	queues map[string]*model.MessageQueue

	// mutex
	mu sync.Mutex
}

func NewCache() *Cache {
	return &Cache{
		rooms:  make(map[string]model.ChatRoom),
		users:  make(map[string]model.User),
		queues: make(map[string]*model.MessageQueue),
	}
}

func (c *Cache) AddNewChatRoom(cr model.ChatRoom) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.rooms[cr.ID] = cr
	return nil
}

func (c *Cache) AddNewUser(user model.User) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	room, ok := c.rooms[user.ChatroomID]
	if !ok {
		return errors.New("chatroom does not exist")
	}

	room.Users = append(room.Users, user.Username)
	c.rooms[user.ChatroomID] = room

	c.queues[user.Username] = &model.MessageQueue{}
	c.users[user.Username] = user

	return nil
}

func (c *Cache) AddMessageToUser(userid string, msg model.Message) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	queue, ok := c.queues[userid]
	if !ok {
		return errors.New("user does not exist")
	}

	queue.Enqueue(msg)
	c.queues[userid] = queue
	return nil
}

func (c *Cache) GetUser(userid string) (model.User, error) {
	user, ok := c.users[userid]
	if !ok {
		return model.User{}, errors.New("user does not exist")
	}

	return user, nil
}

func (c *Cache) GetUsers() ([]model.User, error) {
	result := make([]model.User, 0, len(c.users))

	for _, val := range c.users {
		result = append(result, val)
	}

	return result, nil
}

func (c *Cache) GetUserQueue(userid string) (model.MessageQueue, error) {
	queue, ok := c.queues[userid]
	if !ok {
		return model.MessageQueue{}, errors.New("user does not exist")
	}

	result := model.MessageQueue{}

	for _, val := range *queue {
		result.Enqueue(val)
	}

	return result, nil
}

func (c *Cache) GetRoom(roomid string) (model.ChatRoom, error) {
	room, ok := c.rooms[roomid]
	if !ok {
		return model.ChatRoom{}, errors.New("chatroom does not exist")
	}

	return room, nil
}

func (c *Cache) GetRooms() ([]model.ChatRoom, error) {
	result := make([]model.ChatRoom, 0, len(c.rooms))

	for _, val := range c.rooms {
		result = append(result, val)
	}

	return result, nil
}

func (c *Cache) RemoveUser(userid string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	user, ok := c.users[userid]
	if !ok {
		return errors.New("user does not exist")
	}

	delete(c.users, userid)
	delete(c.queues, userid)

	room, ok := c.rooms[user.ChatroomID]
	if !ok {
		return errors.New("room does not exist")
	}

	index := -1
	for i, val := range room.Users {
		if val == userid {
			index = i
			break
		}
	}

	if index == -1 {
		return nil
	}

	// remove user from room
	room.Users = append(room.Users[:index], room.Users[index+1:]...)
	c.rooms[user.ChatroomID] = room

	return nil
}

func (c *Cache) RemoveRoom(roomid string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.rooms, roomid)

	return nil
}

func (c *Cache) ClearUserQueue(userid string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	queue, ok := c.queues[userid]
	if !ok {
		return errors.New("user does not exist")
	}

	queue.Clear()
	c.queues[userid] = queue

	return nil
}
