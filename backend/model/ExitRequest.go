package model

type ExitRequest struct {
	RoomID   string `json:"chatroom_id"`
	Username string `json:"user_uuid"`
}
