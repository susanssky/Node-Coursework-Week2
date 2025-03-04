package store

import "time"

type Message struct {
	ID       string    `json:"id"`
	From     string    `json:"from"`
	Text     string    `json:"text"`
	TimeSent time.Time `json:"time_sent,omitempty"`
}

type CreateMessageRequest struct {
	From string `json:"from" example:"Alice"`
	Text string `json:"text" example:"Hello World"`
}
