package main

import "github.com/protomem/msg-processor/pkg/ctxstore"

const (
	TraceIDKey = ctxstore.Key("traceId")
	HandlerKey = ctxstore.Key("handler")
)
