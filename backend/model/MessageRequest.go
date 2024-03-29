package model

type MessageRequest struct {
	RoomID   string `json:"chatroom_id"`
	Username string `json:"username"`
	Content  string `json:"message"`
}
