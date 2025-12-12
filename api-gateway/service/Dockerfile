FROM golang:1.22 AS builder

WORKDIR /app
COPY .. .

WORKDIR /app/api-gateway

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /service .

FROM alpine:3.18
WORKDIR /root/
COPY --from=builder /service .

CMD ["./api-gateway"]
