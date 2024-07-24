package main

import (
	"encoding/json"
	"net/http"

	"github.com/protomem/msg-processor/docs"
	"github.com/protomem/msg-processor/pkg/ctxstore"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func (s *APIServer) setupSwagger() {
	docs.SwaggerInfo.Title = "Message Processor"
	docs.SwaggerInfo.Description = "HTTP API - Message Processor"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = s.opts.BaseURL
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http"}
}

func (s *APIServer) setupRoutes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /health", MakeHTTPHandleFunc(s.log, "health", s.handleHealth))

	router.HandleFunc("POST /api/msg", MakeHTTPHandleFunc(s.log, "saveMessage", s.handleSaveMessage))
	router.HandleFunc("GET /api/msg", MakeHTTPHandleFunc(s.log, "messageStatistics", s.handleMessageStatistics))

	router.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL(
			"http://"+s.opts.BaseURL+"/swagger/doc.json",
		),
	))

	return router
}

func (s *APIServer) handleHealth(w http.ResponseWriter, _ *http.Request) error {
	return WriteJSON(w, http.StatusOK, JSONObject{"status": "OK"})
}

// Handle Save Message
//
//	@Summary		Save message
//	@Description	Save message
//	@Tags			message
//	@Accept			json
//	@Produce		json
//	@Param			message	body		SaveMessageDTO	true	"Message"
//	@Success		201		{object}	Message
//	@Failure		500		{object}	any
//	@Router			/msg [post]
func (s *APIServer) handleSaveMessage(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	log := s.log.With(
		HandlerKey.String(), ctxstore.MustFrom[string](ctx, HandlerKey),
		TraceIDKey.String(), ctxstore.MustFrom[string](ctx, TraceIDKey),
	)

	var dto SaveMessageDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return err
	}

	log.Debug("received request")

	msgId, err := s.store.SaveMessage(ctx, dto)
	if err != nil {
		return err
	}

	msg, err := s.store.GetMessage(ctx, msgId)
	if err != nil {
		return err
	}
	msg.Status = MessageProcessing

	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	evt := NewEvent([]byte("newMessage"), msgJSON)
	if err := s.queue.WriteEvents(ctx, evt); err != nil {
		return err
	}

	if err := s.store.UpdateStatusMessages(ctx, []uint64{msgId}, MessageProcessing); err != nil {
		return err
	}

	log.Debug("saved message", "msgId", msg.ID)

	return WriteJSON(w, http.StatusCreated, msg)
}

// Handle Message Statistics
//
//	@Summary		Message statistics
//	@Description	Message statistics
//	@Tags			message
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	MessageStatisticsDTO
//	@Failure		500	{object}	any
//	@Router			/msg [get]
func (s *APIServer) handleMessageStatistics(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	log := s.log.With(
		HandlerKey.String(), ctxstore.MustFrom[string](ctx, HandlerKey),
		TraceIDKey.String(), ctxstore.MustFrom[string](ctx, TraceIDKey),
	)

	var (
		err   error
		stats MessageStatisticsDTO
	)

	stats.Processing, err = s.store.CountProcessingMessages(ctx)
	if err != nil {
		return err
	}

	stats.Completed, err = s.store.CountCompletedMessages(ctx)
	if err != nil {
		return err
	}

	log.Debug("get message statistics")

	return WriteJSON(w, http.StatusOK, stats)
}
