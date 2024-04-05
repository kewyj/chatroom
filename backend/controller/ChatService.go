package controller

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/kewyj/chatroom/model"
	"github.com/kewyj/chatroom/storage"
	"github.com/nwtgck/go-fakelish"
)

const MAX_USERS_IN_ROOM = 10
const MAX_MESSAGES_IN_ROOM = 100

type ChatService struct {
	storage storage.Storage
}

func NewChatService() *ChatService {
	cache := storage.NewCache()

	chatsvc := ChatService{
		storage: cache,
	}

	chatsvc.storage.Initialize()
	return &chatsvc
}

func (cs *ChatService) GetRooms() ([]model.GetRoomsResponse, error) {
	rooms, err := cs.storage.GetRooms()
	if err != nil {
		return nil, err
	}

	var responseSlice []model.GetRoomsResponse

	for _, room := range rooms {
		response := model.GetRoomsResponse{
			RoomID:   room.ID,
			NumUsers: room.UserCount,
		}
		responseSlice = append(responseSlice, response)
	}

	return responseSlice, nil
}

func (cs *ChatService) AddRoom() (string, error) {
	name := fakelish.GenerateFakeWord(6, 9)
	name = strings.Title(name)

	for cs.storage.CheckIfRoomExists(name) {
		name = fakelish.GenerateFakeWord(6, 9)
		name = strings.Title(name)
	}

	err := cs.storage.NewChatRoom(name)
	if err != nil {
		return "", err
	}

	return name, nil
}

func (cs *ChatService) AddUser(user model.NewUserRequest) (string, error) {
	user_uuid := uuid.New().String()

	err := cs.storage.NewUser(user_uuid, user.CustomUsername)
	if err != nil {
		return "", err
	}

	return user_uuid, nil
}

func (cs *ChatService) AddUserToRoom(msg model.AddRoomRequest) error {
	if !cs.storage.CheckIfRoomExists(msg.RoomID) {
		return errors.New("tried to add user to chatroom that does not exist")
	}

	room, err := cs.storage.GetRoom(msg.RoomID)
	if err != nil {
		return err
	}

	if room.UserCount >= MAX_USERS_IN_ROOM {
		return errors.New("chatroom at max capacity")
	}

	username, err := cs.storage.GetUsername(msg.Username)
	if err != nil {
		return err
	}

	cs.sendUserJoinMessage(model.AddRoomRequest{
		RoomID:   msg.RoomID,
		Username: username,
	})

	return cs.storage.AddUserToChatRoom(msg.RoomID)
}

func (cs *ChatService) SendMessage(msg model.MessageRequest) error {
	if !cs.storage.CheckIfRoomExists(msg.RoomID) {
		return errors.New("tried to send message to chatroom that does not exist")
	}

	messages, err := cs.storage.GetRoomMessages(msg.RoomID)
	if err != nil {
		return err
	}

	if len(messages) > MAX_MESSAGES_IN_ROOM {
		cs.storage.RemoveEarliestMessage(msg.RoomID)
	}

	cs.storage.AddMessageToChatRoom(msg.RoomID, model.Message{
		Username: msg.Username,
		Content:  msg.Content,
	})

	//cs.printServerStatus()

	return nil
}

func (cs *ChatService) Poll(req model.PollRequest) ([]model.Message, error) {
	messages, err := cs.storage.GetRoomMessages(req.RoomID)
	if err != nil {
		return []model.Message{}, err
	}

	return messages, nil
}

func (cs *ChatService) RemoveUserFromRoom(req model.ExitRoomRequest) error {
	if !cs.storage.CheckIfRoomExists(req.RoomID) {
		return errors.New("tried to send message to chatroom that does not exist")
	}

	username, err := cs.storage.GetUsername(req.Username)
	if err != nil {
		return err
	}

	err = cs.storage.RemoveUserFromChatRoom(req.RoomID)
	if err != nil {
		return err
	}

	room, err := cs.storage.GetRoom(req.RoomID)
	if err != nil {
		return err
	}

	if room.UserCount == 0 {
		err = cs.storage.RemoveRoom(req.RoomID)
		if err != nil {
			return err
		}
	} else {
		cs.sendUserExitMessage(model.ExitRoomRequest{
			RoomID:   req.RoomID,
			Username: username,
		})
	}

	return nil
}

func (cs *ChatService) RemoveUser(req model.ExitRequest) error {
	return cs.storage.RemoveUser(req.Username)
}

func (cs *ChatService) ClearAll(password string) error {
	if password != "Actually, life is beautiful and I have time." {
		return errors.New("wrong password, database NOT cleared")
	}
	return cs.storage.ClearAll()
}

func (cs *ChatService) Quit(uuid string, chatroom_id string) error {
	err := cs.RemoveUserFromRoom(model.ExitRoomRequest{
		RoomID:   chatroom_id,
		Username: uuid,
	})
	if err != nil {
		return err
	}

	err = cs.RemoveUser(model.ExitRequest{
		Username: uuid,
	})
	if err != nil {
		return err
	}

	return nil
}

func (cs *ChatService) sendUserJoinMessage(user model.AddRoomRequest) {
	cs.storage.AddMessageToChatRoom(user.RoomID, model.Message{
		Username: ".:system:.",
		Content:  user.Username + " has joined the chat.",
	})
}

func (cs *ChatService) sendUserExitMessage(req model.ExitRoomRequest) {
	cs.storage.AddMessageToChatRoom(req.RoomID, model.Message{
		Username: ".:system:.",
		Content:  req.Username + " has left the chat.",
	})
}

func (cs *ChatService) printServerStatus() {
	rooms, _ := cs.storage.GetRooms()

	for _, val := range rooms {
		fmt.Println("ROOM: " + val.ID)

		fmt.Println("\tMessages:")
		for _, msg := range val.Messages {
			fmt.Println("\t\t" + msg.Username)
			fmt.Println("\t\t\t" + msg.Content)
		}
	}
}
