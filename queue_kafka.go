package main

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/protomem/msg-processor/pkg/ctxstore"
	"github.com/segmentio/kafka-go"
)

var _ Queue = (*KafkaQueue)(nil)

type KafkaQueueOptions struct {
	Addrs string
	Topic string
}

type KafkaQueue struct {
	opts KafkaQueueOptions
	log  *slog.Logger

	writer *kafka.Writer
	reader *kafka.Reader
}

func NewKafkaQueue(ctx context.Context, log *slog.Logger, opts KafkaQueueOptions) (*KafkaQueue, error) {
	addrs := strings.Split(opts.Addrs, ",")

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(addrs...),
		Topic:                  opts.Topic,
		Balancer:               &kafka.Hash{},
		AllowAutoTopicCreation: true,
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   addrs,
		Topic:     opts.Topic,
		GroupID:   "msg-processor",
		Partition: 0,
		MaxBytes:  10e6, // 10MB
	})

	return &KafkaQueue{
		opts: opts,
		log:  log.With("component", "kafkaQueue"),

		writer: writer,
		reader: reader,
	}, nil
}

func (q *KafkaQueue) WriteEvents(ctx context.Context, events ...Event) error {
	log := q.log.With(TraceIDKey.String(), ctxstore.MustFrom[string](ctx, TraceIDKey))
	msgs := kafakMsgsFromEvents(events...)

	if err := q.writer.WriteMessages(ctx, msgs...); err != nil {
		log.Debug("failed to write events", "error", err)

		return err
	}

	log.Debug("written events", "countEvents", len(msgs))

	return nil
}

func (q *KafkaQueue) ReadEvent(ctx context.Context) (Event, error) {
	log := q.log.With(TraceIDKey.String(), ctxstore.MustFrom[string](ctx, TraceIDKey))

	msg, err := q.reader.ReadMessage(ctx)
	if err != nil {
		log.Debug("failed to read event", "error", err)

		return Event{}, err
	}

	log.Debug("read event")

	return eventFromKafkaMsg(msg), nil
}

func (q *KafkaQueue) Close(_ context.Context) error {
	var errs error

	if err := q.writer.Close(); err != nil {
		errs = errors.Join(errs, err)
	}

	if err := q.reader.Close(); err != nil {
		errs = errors.Join(errs, err)
	}

	return errs
}

func kafkaMsgFromEvent(evt Event) kafka.Message {
	return kafka.Message{
		Key:   bytes.Clone(evt.Key),
		Value: bytes.Clone(evt.Value),
	}
}

func eventFromKafkaMsg(msg kafka.Message) Event {
	return Event{
		Key:    bytes.Clone(msg.Key),
		Value:  bytes.Clone(msg.Value),
		Tstamp: msg.Time,
	}
}

func kafakMsgsFromEvents(evts ...Event) []kafka.Message {
	msgs := make([]kafka.Message, len(evts))
	for i := 0; i < len(evts); i++ {
		msgs[i] = kafkaMsgFromEvent(evts[i])
	}
	return msgs
}
