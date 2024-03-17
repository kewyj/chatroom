package model

type MessageRequest struct {
	RoomID   string `json:"chatroom_id"`
	Username string `json:"user_uuid"`
	Content  string `json:"message"`
}
