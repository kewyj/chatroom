package model

type PollRequest struct {
	RoomID   string `json:"chatroom_id"`
	Username string `json:"user_uuid"`
}
