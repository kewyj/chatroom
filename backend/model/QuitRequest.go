package model

type QuitRequest struct {
	RoomID   string `json:"chatroom_id"`
	Username string `json:"user_uuid"`
}
