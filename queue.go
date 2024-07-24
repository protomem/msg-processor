package main

import (
	"context"
	"time"
)

type Event struct {
	Key    []byte
	Value  []byte
	Tstamp time.Time
}

func NewEvent(key []byte, value []byte) Event {
	return Event{
		Key:    key,
		Value:  value,
		Tstamp: time.Now(),
	}
}

type Queue interface {
	WriteEvents(ctx context.Context, events ...Event) error
	ReadEvent(ctx context.Context) (Event, error)

	Close(ctx context.Context) error
}
