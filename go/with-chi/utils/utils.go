package utils

import (
	"encoding/json"
	"net/http"
	"node-week-02-with-chi/store"
)

// CreateMessage and GetMessage return store.Message
// UpdateMessage returns *store.Message
// GetAllMessages, GetLatestMessages, and GetSearchedMessages return []store.Message
type MessageData interface {
	store.Message | *store.Message | []store.Message
}

type Response[T MessageData] struct {
	Data T `json:"data"`
}
type ErrorResponse struct {
	Error string `json:"error"`
}

func ParseJSON(r *http.Request, payload any) error {
	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON[T MessageData](w http.ResponseWriter, status int, data T) error {
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(Response[T]{Data: data})
}

func WriteError(w http.ResponseWriter, status int, err string) error {
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(ErrorResponse{Error: err})
}
