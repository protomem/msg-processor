package main

import (
	"context"
	"log/slog"
)

var _ Queue = (*KafkaQueue)(nil)

type KafkaQueueOptions struct {
	Addrs string
	Topic string
}

type KafkaQueue struct {
	log *slog.Logger
}

func NewKafkaQueue(ctx context.Context, log *slog.Logger, opts KafkaQueueOptions) (*KafkaQueue, error) {
	return &KafkaQueue{}, nil
}

func (q *KafkaQueue) WriteEvents(ctx context.Context, events ...Event) error {
	panic("unimplemented")
}

func (q *KafkaQueue) Close(_ context.Context) error {
	panic("unimplemented")
}
