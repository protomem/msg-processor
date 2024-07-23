package main

import (
	"encoding/json"
	"net/http"
)

func (s *APIServer) setupRoutes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /health", MakeHTTPHandleFunc(s.handleHealth))

	router.HandleFunc("POST /api/messages", MakeHTTPHandleFunc(s.handleSaveMessage))

	return router
}

func (s *APIServer) handleHealth(w http.ResponseWriter, _ *http.Request) error {
	return WriteJSON(w, http.StatusOK, JSONObject{"status": "OK"})
}

func (s *APIServer) handleSaveMessage(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	log := s.log.With("handler", "saveMessage")

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

	// TODO: Send message to queue(kafka)

	if err := s.store.UpdateStatusMessages(ctx, []uint64{msgId}, MessageProcessing); err != nil {
		return err
	}

	log.Debug("saved message", "msgId", msg.ID)

	return WriteJSON(w, http.StatusOK, msg)
}
