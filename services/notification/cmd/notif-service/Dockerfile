FROM golang:1.22 AS builder

WORKDIR /app
COPY ../.. .

WORKDIR /app/services/notification/cmd/notification-service

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /notif-service .

FROM alpine:3.18
WORKDIR /root/

COPY --from=builder /notif-service .

CMD ["./notification-service"]
