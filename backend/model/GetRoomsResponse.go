package model

type GetRoomsResponse struct {
	RoomID   string `json:"chatroom_id"`
	NumUsers int    `json:"num_users"`
}
