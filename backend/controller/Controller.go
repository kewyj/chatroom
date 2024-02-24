package controller

import "github.com/kewyj/chatroom/model"

type Controller interface {
	AddUser() (string, error)
	SendMessage(msg model.Message) error
	Poll(string) (model.MessageQueue, error)
	RemoveUser(string) error
	IsUserSpamming(string) bool
}
