FROM golang:1.22 AS builder

WORKDIR /app
COPY ../.. .

WORKDIR /app/services/post/cmd/post-service

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /post-service .

FROM alpine:3.18
WORKDIR /root/

COPY --from=builder /post-service .

CMD ["./post-service"]
