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

// CreateMessage godoc
// @Summary Create a message
// @Description Create a new message and add it to the system
// @Tags messages
// @Accept json
// @Produce json
// @Param message body store.CreateMessageRequest true "Message content"
// @Success 201 {object} utils.Response[store.Message] "Successful creation of message"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /messages [post]
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

// GetAllMessages godoc
// @Summary Get all messages
// @Description Return a list of all messages in the app
// @Tags messages
// @Produce json
// @Success 200 {object} utils.Response[[]store.Message] "message list"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /messages [get]
func (h *MessageHandler) GetAllMessages(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, h.Message)
}

// GetLatestMessages godoc
// @Summary Get the latest 10 messages
// @Description Return the latest 10 messages
// @Tags messages
// @Produce json
// @Success 200 {object} utils.Response[[]store.Message] "the latest 10 messages list"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /messages/latest [get]
func (h *MessageHandler) GetLatestMessages(w http.ResponseWriter, r *http.Request) {
	start := max(0, len(h.Message)-10)
	respondJSON(w, http.StatusOK, h.Message[start:])
}

// GetSearchedMessages godoc
// @Summary Get the messages that has searched if matched
// @Description Return the messages that has searched if matched
// @Tags messages
// @Produce json
// @Param text query string true "Text to search for in messages"
// @Success 200 {object} utils.Response[[]store.Message] "the messages list if matched"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 404 {object} utils.ErrorResponse "No matching messages found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /messages/search [get]
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

// GetMessage godoc
// @Summary Get a message by ID
// @Description Return a message by ID
// @Tags messages
// @Produce json
// @Param messageId path string true "Message ID"
// @Success 200 {object} utils.Response[store.Message] "the message if matched"
// @Failure 404 {object} utils.ErrorResponse "No matching messages found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /messages/{messageId} [get]
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

// UpdateMessage godoc
// @Summary Update a message by ID
// @Description Return an updated message by ID
// @Tags messages
// @Accept json
// @Produce json
// @Param messageId path string true "Message ID"
// @Param message body store.CreateMessageRequest true "Updated message content"
// @Success 200 {object} utils.Response[store.Message] "the updated message"
// @Failure 404 {object} utils.ErrorResponse "No matching messages found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /messages/{messageId} [put]
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

// DeleteMessage godoc
// @Summary Delete a message by ID
// @Description Deletes a message with the specified ID and returns no content on success
// @Tags messages
// @Produce json
// @Param messageId path string true "Message ID"
// @Success 204 "No Content - Message successfully deleted"
// @Failure 404 {object} utils.ErrorResponse "No matching message found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /messages/{messageId} [delete]
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
