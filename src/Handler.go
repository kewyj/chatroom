package src

import (
	"encoding/json"
	"errors"
	"io"
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
	msg, err := unmarshalMessage(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	err = h.controller.SendMessage(msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Poll(w http.ResponseWriter, r *http.Request) {
	msg, err := unmarshalMessage(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	messages := h.controller.RetrieveUndelivered(msg.Username)

	// Convert the array of messages to JSON format
	jsonResponse, err := json.Marshal(*messages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	messages.Clear()
	// Set the content type header
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON-encoded array to the response body
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func unmarshalMessage(rbody io.ReadCloser) (Message, error) {
	body, err := ioutil.ReadAll(rbody)
	if err != nil {
		return Message{}, errors.New("error reading request body")
	}

	var msg Message
	if err = json.Unmarshal(body, &msg); err != nil {
		return Message{}, errors.New("error decoding JSON")
	}

	return msg, nil
}
