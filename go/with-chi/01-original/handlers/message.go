package handlers

import (
	"net/http"
	"node-week-02-with-chi/store"
	"node-week-02-with-chi/utils"

	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

type MessageHandler struct {
	Message []store.Message
}

func New() *MessageHandler {
	return &MessageHandler{
		Message: messages,
	}
}

var messages = []store.Message{
	{
		ID:   strconv.Itoa(0),
		From: "Bart",
		Text: "Welcome to CYF chat system!",
	},
}

func validateMessage(req store.CreateMessageRequest) bool {
	return req.From != "" && req.Text != ""
}

func (h *MessageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var req store.CreateMessageRequest

	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !validateMessage(req) {
		utils.WriteError(w, http.StatusBadRequest, "Your name or message are missing.")
		return
	}

	newMessage := store.Message{
		ID:       strconv.Itoa(len(h.Message)),
		From:     req.From,
		Text:     req.Text,
		TimeSent: time.Now().UTC(),
	}

	h.Message = append(h.Message, newMessage)

	respondJSON(w, http.StatusCreated, newMessage)

}

func (h *MessageHandler) GetAllMessages(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, h.Message)
}

func (h *MessageHandler) GetLatestMessages(w http.ResponseWriter, r *http.Request) {
	start := max(0, len(h.Message)-10)
	latestMessages := h.Message[start:]

	reversed := make([]store.Message, len(latestMessages))
	for i, j := 0, len(latestMessages)-1; j >= 0; i, j = i+1, j-1 {
		reversed[i] = latestMessages[j]
	}

	respondJSON(w, http.StatusOK, reversed)
}

func (h *MessageHandler) GetSearchedMessages(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")

	if text == "" {
		utils.WriteError(w, http.StatusBadRequest, "Please fill the text field")
		return
	}

	var matchedMessages []store.Message
	for _, message := range h.Message {
		if strings.Contains(strings.ToLower(message.Text), strings.ToLower(text)) {
			matchedMessages = append(matchedMessages, message)
		}
	}
	if len(matchedMessages) == 0 {
		utils.WriteError(w, http.StatusNotFound, "Not found the message that has matched")
		return
	}

	respondJSON(w, http.StatusOK, matchedMessages)
}

func (h *MessageHandler) GetMessage(w http.ResponseWriter, r *http.Request) {
	messageId := chi.URLParam(r, "messageId")
	for _, message := range h.Message {
		if message.ID == messageId {
			respondJSON(w, http.StatusOK, message)
			return
		}
	}
	utils.WriteError(w, http.StatusNotFound, "Message not found")
}

func (h *MessageHandler) UpdateMessage(w http.ResponseWriter, r *http.Request) {

	messageId := chi.URLParam(r, "messageId")

	var existingMessage *store.Message
	var req store.CreateMessageRequest

	for i := range h.Message {
		if h.Message[i].ID == messageId {
			existingMessage = &h.Message[i]
			break
		}
	}

	if existingMessage == nil {
		utils.WriteError(w, http.StatusNotFound, "Message not found")
		return
	}

	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !validateMessage(req) {
		utils.WriteError(w, http.StatusBadRequest, "Your name or message are missing.")
		return
	}

	existingMessage.From = req.From
	existingMessage.Text = req.Text

	respondJSON(w, http.StatusOK, existingMessage)
}

func (h *MessageHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	messageId := chi.URLParam(r, "messageId")

	index := -1
	for i, message := range h.Message {
		if message.ID == messageId {
			index = i
			break
		}
	}

	if index == -1 {
		utils.WriteError(w, http.StatusNotFound, "Message not found")
		return
	}

	h.Message = slices.Delete(h.Message, index, index+1)

	w.WriteHeader(http.StatusNoContent)
}

func respondJSON[T utils.MessageData](w http.ResponseWriter, status int, data T) {
	if err := utils.WriteJSON(w, status, data); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
	}
}
