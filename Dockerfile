FROM golang:1.22.4-alpine AS builder

WORKDIR /app

RUN go install std

COPY go.* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/msg-processor .



FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/msg-processor .
COPY --from=builder /app/.env .

CMD ["./msg-processor", "-cfg", ".env"]