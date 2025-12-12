FROM golang:1.22 AS builder

WORKDIR /app
COPY ../.. .

WORKDIR /app/services/search/cmd/search-service

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /search-service .

FROM alpine:3.18
WORKDIR /root/

COPY --from=builder /search-service .

CMD ["./search-service"]
