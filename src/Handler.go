package src

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type Handler struct {
	controller Controller
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
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   name,
		MaxAge: -1,
	})
	w.WriteHeader(http.StatusOK)
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