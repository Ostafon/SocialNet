FROM golang:1.22 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o chat-service ./cmd/chat-service/

FROM alpine:3.18
WORKDIR /root/
COPY --from=builder /app/chat-service .

EXPOSE 8080
EXPOSE 50056 8082
CMD ["./chat-service"]
