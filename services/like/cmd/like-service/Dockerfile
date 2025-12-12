FROM golang:1.22 AS builder

WORKDIR /app
COPY ../.. .

WORKDIR /app/services/like/cmd/like-service

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /like-service .

FROM alpine:3.18
WORKDIR /root/

COPY --from=builder /like-service .

CMD ["./like-service"]
