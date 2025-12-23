FROM golang:1.22 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o gateway ./cmd/gateway/

FROM alpine:3.18
WORKDIR /root/
COPY --from=builder /app/gateway .

EXPOSE 8080
CMD ["./gateway"]
