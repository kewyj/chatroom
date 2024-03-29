package model

type NewUserRequest struct {
	CustomUsername string `json:"custom_username"`
	RoomID         string `json:"chatroom_id"`
}
