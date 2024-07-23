package main

import (
	"context"
	"database/sql"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const _pgDriverName = "pgx"

var _ Storage = (*PgStorage)(nil)

type PgStorageOptions struct {
	DSN  string
	Ping bool
}

type PgStorage struct {
	opts PgStorageOptions
	log  *slog.Logger
	db   *sql.DB
}

func NewPgStorage(ctx context.Context, log *slog.Logger, opts PgStorageOptions) (*PgStorage, error) {
	db, err := sql.Open(_pgDriverName, opts.DSN)
	if err != nil {
		return nil, err
	}

	if opts.Ping {
		if err := db.PingContext(ctx); err != nil {
			return nil, err
		}
	}

	return &PgStorage{
		opts: opts,
		log:  log.With("component", "pgStorage"),
	}, nil
}

func (s *PgStorage) Close(ctx context.Context) error {
	return s.db.Close()
}

func (s *PgStorage) GetMessage(ctx context.Context, id uint64) (msg Message, err error) {
	panic("unimplemented")
}

func (s *PgStorage) SaveMessage(ctx context.Context, msg Message) (id uint64, err error) {
	panic("unimplemented")
}

func (s *PgStorage) UpdateStatusMessages(ctx context.Context, ids []uint64, status MessageStatus) error {
	panic("unimplemented")
}
