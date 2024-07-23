package main

import (
	"context"
	"database/sql"
	"errors"
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

func (s *PgStorage) GetMessage(ctx context.Context, id uint64) (Message, error) {
	log := s.log.With("query", "getMessage")

	query := `
		SELECT (id, created_at, updated_at, message, status)
		FROM messages
		WHERE id = $1
		LIMIT 1
	`

	log.Debug("build query", "query", query, "args", []any{id})

	var msg Message
	row := s.db.QueryRowContext(ctx, query, id)
	if err := row.Scan(&msg.ID, &msg.CreatedAt, &msg.UpdatedAt, &msg.Text, &msg.Status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Message{}, ErrMsgNotFound
		}

		log.Debug("failed to execute query", "error", err)

		return Message{}, err
	}

	log.Debug("executed query", "result", msg)

	return msg, nil
}

func (s *PgStorage) SaveMessage(ctx context.Context, dto SaveMessageDTO) (uint64, error) {
	log := s.log.With("query", "saveMessage")

	query := `
		INSERT INTO messages (message)
		VALUES ($1)
		RETURNING id
	`

	log.Debug("build query", "query", query, "args", []any{dto.Text})

	var id uint64
	row := s.db.QueryRowContext(ctx, query, dto.Text)
	if err := row.Scan(&id); err != nil {
		log.Debug("failed to execute query", "error", err)

		return 0, err
	}

	log.Debug("executed query", "result", id)

	return id, nil
}

func (s *PgStorage) UpdateStatusMessages(ctx context.Context, ids []uint64, status MessageStatus) error {
	log := s.log.With("query", "updateStatusMessages")

	query := `
		UPDATE messages
		SET status = $1
		WHERE id = ANY($2::bigint[])
	`

	log.Debug("build query", "query", query, "args", []any{status, ids})

	res, err := s.db.ExecContext(ctx, query, status, ids)
	if err != nil {
		log.Debug("failed to execute query", "error", err)

		return err
	}

	countRows, err := res.RowsAffected()
	if err != nil {
		log.Debug("failed to get rows count", "error", err)
	}

	log.Debug("executed query", "updatedRows", countRows)

	return nil
}
