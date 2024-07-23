package main

import "time"

type MessageStatus string

const (
	MessageCreated    MessageStatus = "created"
	MessageProcessing MessageStatus = "processing"
	MessageCompleted  MessageStatus = "completed"
)

type Message struct {
	ID uint64 `json:"id"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Text string `json:"text"`

	Status MessageStatus `json:"status"`
}
