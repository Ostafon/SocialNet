FROM golang:1.22 AS builder

WORKDIR /app
COPY ../.. .

WORKDIR /app/services/comment/cmd/comment-service

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /comment-service .

FROM alpine:3.18
WORKDIR /root/

COPY --from=builder /comment-service .

CMD ["./comment-service"]
