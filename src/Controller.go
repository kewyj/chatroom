package src

type Controller interface {
	AddUser() (string, error)
	SendMessage(msg Message) error
}