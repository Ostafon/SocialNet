FROM golang:1.22 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o post-service ./cmd/post-service/

FROM alpine:3.18
WORKDIR /root/
COPY --from=builder /app/post-service .

EXPOSE 8080
EXPOSE 50053
CMD ["./post-service"]
