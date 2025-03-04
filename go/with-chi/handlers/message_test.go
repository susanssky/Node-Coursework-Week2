package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"node-week-02-with-chi/store"
	"node-week-02-with-chi/utils"

	"github.com/go-chi/chi/v5"
)

// TestNew tests whether the New() function initialises the MessageHandler correctly
func TestNew(t *testing.T) {
	handler := New()

	t.Run("Check if the returned handler is nil", func(t *testing.T) {
		if handler == nil {
			t.Errorf("New() returned nil, expected a non-null MessageHandler")
			return
		}
	})

	t.Run("Check if the Message field is initialised correctly", func(t *testing.T) {
		if handler.Message == nil {
			t.Errorf("The Message field is nil and is expected to contain the initial message")
		}
	})

	t.Run("Checking the length of the Message slice", func(t *testing.T) {
		expectedLength := len(messages)
		if len(handler.Message) != expectedLength {
			t.Errorf("Expected Message slice length %d, actually got %d", expectedLength, len(handler.Message))
		}
	})

	t.Run("Check if the initial message content is correct", func(t *testing.T) {
		if len(handler.Message) > 0 {
			firstMessage := handler.Message[0]
			if firstMessage.ID != "0" || firstMessage.From != "Bart" || firstMessage.Text != "Welcome to CYF chat system!" {
				t.Errorf("The initial message content is incorrect, got %+v", firstMessage)
			}
		}
	})
}

// Initialisation function for testing to ensure consistent test environment
func setupTestHandler() *MessageHandler {
	return &MessageHandler{
		Message: []store.Message{
			{ID: "0", From: "Bart", Text: "Welcome to CYF chat system!", TimeSent: time.Now().UTC()},
			{ID: "1", From: "Lisa", Text: "Hello everyone!", TimeSent: time.Now().UTC()},
		},
	}
}

// Testing CreateMessage
func TestCreateMessage(t *testing.T) {
	handler := setupTestHandler()

	t.Run("Create message normally", func(t *testing.T) {
		const (
			testName = "Tom"
			testText = "Hello"
		)

		newMessage := store.CreateMessageRequest{From: testName, Text: testText} // 使用 CreateMessageRequest
		body, _ := json.Marshal(newMessage)
		req, _ := http.NewRequest("POST", "/api/v1/messages", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		handler.CreateMessage(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("Expected status code %v, got %v", http.StatusCreated, status)
		}

		var response utils.Response[store.Message]
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}
		createdMessage := response.Data
		if createdMessage.From != testName || createdMessage.Text != testText {
			t.Errorf("Incorrect message content: %+v", createdMessage)
		}
	})

	t.Run("Missing From of Required Fields", func(t *testing.T) {
		invalidMessage := store.CreateMessageRequest{From: "", Text: "Invalid"}
		body, _ := json.Marshal(invalidMessage)
		req, _ := http.NewRequest("POST", "/api/v1/messages", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		handler.CreateMessage(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Expected status code %v, got %v", http.StatusBadRequest, status)
		}
	})

	t.Run("Missing Text of Required Fields", func(t *testing.T) {
		invalidMessage := store.CreateMessageRequest{From: "Invalid", Text: ""}
		body, _ := json.Marshal(invalidMessage)
		req, _ := http.NewRequest("POST", "/api/v1/messages", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		handler.CreateMessage(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Expected status code %v, got %v", http.StatusBadRequest, status)
		}
	})
}

// Testing GetAllMessages
func TestGetAllMessages(t *testing.T) {
	handler := setupTestHandler()

	t.Run("Get All Messages", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/messages", nil)
		rr := httptest.NewRecorder()

		handler.GetAllMessages(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Expected status code %v, got %v", http.StatusOK, status)
		}

		var response utils.Response[[]store.Message]
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}
		messages := response.Data
		if len(messages) != 2 {
			t.Errorf("Expected 2 messages, got %v", len(messages))
		}
	})
}

// Testing GetLatestMessages
func TestGetLatestMessages(t *testing.T) {
	handler := setupTestHandler()

	t.Run("Get Latest 10 Messages", func(t *testing.T) {
		// Add more messages to test "latest 10" logic
		for i := 2; i < 15; i++ {
			handler.Message = append(handler.Message, store.Message{
				ID:   strconv.Itoa(i),
				From: "Test",
				Text: "Message " + strconv.Itoa(i),
			})
		}

		req, _ := http.NewRequest("GET", "/api/v1/messages/latest", nil)
		rr := httptest.NewRecorder()

		handler.GetLatestMessages(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Expected status code %v, got %v", http.StatusOK, status)
		}

		var response utils.Response[[]store.Message]
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}
		messages := response.Data
		if len(messages) != 10 {
			t.Errorf("Expected 10 messages, got %v", len(messages))
		}
	})
}

// Testing GetSearchedMessages
func TestGetSearchedMessages(t *testing.T) {
	handler := setupTestHandler()

	t.Run("Searching for existing message", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/messages/search?text=Welcome", nil)
		rr := httptest.NewRecorder()

		handler.GetSearchedMessages(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Expected status code %v, got %v", http.StatusOK, status)
		}

		var response utils.Response[[]store.Message]
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}
		messages := response.Data
		if len(messages) != 1 {
			t.Errorf("Expected 1 message found, got %v", len(messages))
		}
	})

	t.Run("Search for non-existent messages", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/messages/search?text=Nonexistent", nil)
		rr := httptest.NewRecorder()

		handler.GetSearchedMessages(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("Expected status code %v, got %v", http.StatusNotFound, status)
		}
	})
}

// Testing GetMessage
func TestGetMessage(t *testing.T) {
	handler := setupTestHandler()

	t.Run("Get a message by id", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/messages/0", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("messageId", "0")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		rr := httptest.NewRecorder()

		handler.GetMessage(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Expected status code %v, got %v", http.StatusOK, status)
		}

		var response utils.Response[store.Message]
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}
		message := response.Data
		if message.ID != "0" {
			t.Errorf("Message lookup failed, expected id '0', got: %v", message.ID)
		}
	})
}

// Testing UpdateMessage
func TestUpdateMessage(t *testing.T) {
	handler := setupTestHandler()
	const (
		updatedName = "Marge"
		updatedText = "Updated message"
	)

	t.Run("update a message by id", func(t *testing.T) {
		updatedMessage := store.CreateMessageRequest{From: updatedName, Text: updatedText}
		body, _ := json.Marshal(updatedMessage)
		req, _ := http.NewRequest("PUT", "/api/v1/messages/0", bytes.NewBuffer(body))
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("messageId", "0")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		rr := httptest.NewRecorder()

		handler.UpdateMessage(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Expected status code %v, got %v", http.StatusOK, status)
		}

		var response utils.Response[*store.Message]
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}
		updated := response.Data
		if updated.From != updatedName || updated.Text != updatedText {
			t.Errorf("Message update failed: %+v", updated)
		}
	})
}

// Testing DeleteMessage
func TestDeleteMessage(t *testing.T) {
	handler := setupTestHandler()

	t.Run("delete a message by id", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/api/v1/messages/0", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("messageId", "0")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		rr := httptest.NewRecorder()

		handler.DeleteMessage(rr, req)

		if status := rr.Code; status != http.StatusNoContent {
			t.Errorf("Expected status code %v, got %v", http.StatusNoContent, status)
		}

		if len(handler.Message) != 1 {
			t.Errorf("Expected 1 message, got %v", len(handler.Message))
		}
	})
}
