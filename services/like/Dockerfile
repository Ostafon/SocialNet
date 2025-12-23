FROM golang:1.22 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o like-service ./cmd/like-service/

FROM alpine:3.18
WORKDIR /root/
COPY --from=builder /app/like-service .

EXPOSE 8080
EXPOSE 50055
CMD ["./like-service"]
