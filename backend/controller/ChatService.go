package controller

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/kewyj/chatroom/model"
	"github.com/kewyj/chatroom/storage"
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
			NumUsers: len(room.Users),
		}
		responseSlice = append(responseSlice, response)
	}

	return responseSlice, nil
}

func (cs *ChatService) AddRoom() (string, error) {
	newroom := uuid.New().String()

	err := cs.storage.NewChatRoom(newroom)
	if err != nil {
		return "", err
	}

	return newroom, nil
}

func (cs *ChatService) AddUser(user model.NewUserRequest) (string, error) {
	if !cs.storage.CheckIfRoomExists(user.RoomID) {
		return "", errors.New("tried to add user to chatroom that does not exist")
	}

	users, err := cs.storage.GetRoomUserUUIDs(user.RoomID)
	if err != nil {
		return "", err
	}

	if len(users) >= MAX_USERS_IN_ROOM {
		return "", errors.New("chatroom at max capacity")
	}

	// chatroom exists and can accomodate new user
	user_uuid := uuid.New().String()
	if err := cs.storage.AddUserToChatRoom(user.CustomUsername, user_uuid, user.RoomID); err != nil {
		return "", err
	}

	// user successfully added
	cs.sendUserJoinMessage(user)

	//cs.printServerStatus()

	return user_uuid, nil
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

	username, err := cs.storage.GetUsername(msg.RoomID, msg.Username)
	if err != nil {
		return err
	}

	cs.storage.AddMessageToChatRoom(msg.RoomID, model.Message{
		Username: username,
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

func (cs *ChatService) RemoveUser(req model.ExitRequest) error {
	if !cs.storage.CheckIfRoomExists(req.RoomID) {
		return errors.New("tried to send message to chatroom that does not exist")
	}

	user, err := cs.storage.GetUsername(req.RoomID, req.Username)
	if err != nil {
		return err
	}

	cs.storage.RemoveUserFromChatRoom(req.Username, req.RoomID)

	users, err := cs.storage.GetRoomUsernames(req.RoomID)
	if err != nil {
		return err
	}

	if len(users) == 0 {
		// room empty, delete
		cs.storage.RemoveRoom(req.RoomID)
	} else {
		cs.sendUserExitMessage(model.ExitRequest{
			RoomID:   req.RoomID,
			Username: user,
		})
	}

	//cs.printServerStatus()

	return nil
}

func (cs *ChatService) ClearAll(password string) error {
	if password != "Actually, life is beautiful and I have time." {
		return errors.New("wrong password, database NOT cleared")
	}
	return cs.storage.ClearAll()
}

func (cs *ChatService) sendUserJoinMessage(user model.NewUserRequest) {
	cs.storage.AddMessageToChatRoom(user.RoomID, model.Message{
		Username: ".:system:.",
		Content:  user.CustomUsername + " has joined the chat.",
	})
}

func (cs *ChatService) sendUserExitMessage(user model.ExitRequest) {
	cs.storage.AddMessageToChatRoom(user.RoomID, model.Message{
		Username: ".:system:.",
		Content:  user.Username + " has left the chat.",
	})
}

func (cs *ChatService) printServerStatus() {
	rooms, _ := cs.storage.GetRooms()

	for _, val := range rooms {
		fmt.Println("ROOM: " + val.ID)

		fmt.Println("\tUsers:")
		for _, user := range val.Users {
			fmt.Println("\t\t" + user)
		}

		fmt.Println("\tMessages:")
		for _, msg := range val.Messages {
			fmt.Println("\t\t" + msg.Username)
			fmt.Println("\t\t\t" + msg.Content)
		}
	}
}
