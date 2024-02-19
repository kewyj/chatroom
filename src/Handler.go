package src

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Handler struct {
	controller Controller
}

type NewUserResponse struct {
	Username string `json:"username"`
}

func NewHandler(c Controller) *Handler {
	return &Handler{
		controller: c,
	}
}

func (h *Handler) NewUser(w http.ResponseWriter, r *http.Request) {
	name, err := h.controller.AddUser()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	// return username to be stored as cookie
	response := NewUserResponse{
		Username: name,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var msg Message
	if err = json.Unmarshal(body, &msg); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	err = h.controller.SendMessage(msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
}
