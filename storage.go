package main

import (
	"context"
)

type SaveMessageDTO struct {
	Text string
}

type Storage interface {
	GetMessage(ctx context.Context, id uint64) (msg Message, err error)
	SaveMessage(ctx context.Context, msg Message) (id uint64, err error)
	UpdateStatusMessages(ctx context.Context, ids []uint64, status MessageStatus) error

	Close(ctx context.Context) error
}
