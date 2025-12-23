FROM golang:1.22 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o notification-service ./cmd/notification-service/

FROM alpine:3.18
WORKDIR /root/
COPY --from=builder /app/notification-service .

EXPOSE 8080
EXPOSE 50057 8082
CMD ["./notification-service"]
