FROM golang:1.22 AS builder

WORKDIR /app
COPY ../.. .

WORKDIR /app/services/chat/cmd/chat-service

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /chat-service .

FROM alpine:3.18
WORKDIR /root/

COPY --from=builder /chat-service .

CMD ["./chat-service"]
