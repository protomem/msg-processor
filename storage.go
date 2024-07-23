package main

import (
	"context"
)

type Storage interface {
	GetMessage(ctx context.Context, id uint64) (msg Message, err error)
	SaveMessage(ctx context.Context, dto SaveMessageDTO) (id uint64, err error)
	UpdateStatusMessages(ctx context.Context, ids []uint64, status MessageStatus) error

	Close(ctx context.Context) error
}
