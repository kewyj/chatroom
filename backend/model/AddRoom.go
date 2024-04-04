package model

type AddRoomRequest struct {
	RoomID   string `json:"chatroom_id"`
	Username string `json:"user_uuid"`
}
