package handler

// import (
// 	"encoding/json"
// 	"errors"
// 	"io"
// 	"io/ioutil"
// 	"net/http"

// 	"github.com/kewyj/chatroom/controller"
// 	"github.com/kewyj/chatroom/model"
// )

// type Handler struct {
// 	controller controller.Controller
// }

// func NewHandler(c controller.Controller) *Handler {
// 	return &Handler{
// 		controller: c,
// 	}
// }

// func (h *Handler) NewUser(w http.ResponseWriter, r *http.Request) {
// 	name, err := h.controller.AddUser()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusServiceUnavailable)
// 		return
// 	}
// 	// return username to be stored as cookie
// 	response := NewUserResponse{
// 		Username: name,
// 	}
// 	responseJSON, err := json.Marshal(response)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(responseJSON)
// }

// func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
// 	msg, err := unmarshalMessage(r.Body)
// 	if err != nil {
// 		http.Error(w, "Error reading request body", http.StatusInternalServerError)
// 		return
// 	}

// 	if h.controller.IsUserSpamming(msg.Username) {
// 		response := NewUserResponse{
// 			Username: "$",
// 		}

// 		responseJSON, _ := json.Marshal(response)
// 		w.Header().Set("Content-Type", "application/json")
// 		w.Write(responseJSON)
// 		return
// 	}

// 	err = h.controller.SendMessage(msg)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusServiceUnavailable)
// 		return
// 	}
// 	w.WriteHeader(http.StatusOK)
// }

// func (h *Handler) Poll(w http.ResponseWriter, r *http.Request) {
// 	msg, err := unmarshalMessage(r.Body)
// 	if err != nil {
// 		http.Error(w, "Error reading request body", http.StatusInternalServerError)
// 		return
// 	}
// 	messages, err := h.controller.Poll(msg.Username)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Convert the array of messages to JSON format
// 	jsonResponse, err := json.Marshal(messages)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Set the content type header
// 	w.Header().Set("Content-Type", "application/json")

// 	// Write the JSON-encoded array to the response body
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(jsonResponse)
// }

// func (h *Handler) Exit(w http.ResponseWriter, r *http.Request) {
// 	msg, err := unmarshalMessage(r.Body)
// 	if err != nil {
// 		http.Error(w, "Error reading request body", http.StatusInternalServerError)
// 		return
// 	}
// 	err = h.controller.RemoveUser(msg.Username)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.WriteHeader(http.StatusOK)
// }

// func unmarshalMessage(rbody io.ReadCloser) (model.Message, error) {
// 	body, err := ioutil.ReadAll(rbody)
// 	if err != nil {
// 		return model.Message{}, errors.New("error reading request body")
// 	}

// 	var msg model.Message
// 	if err = json.Unmarshal(body, &msg); err != nil {
// 		return model.Message{}, errors.New("error decoding JSON")
// 	}

// 	return msg, nil
// }

// func enableCORS(w *http.ResponseWriter) {
// 	(*w).Header().Set("Access-Control-Allow-Origin", "*")
// }
