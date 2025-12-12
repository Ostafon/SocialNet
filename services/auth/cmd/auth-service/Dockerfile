# ---------- BUILDER ----------
FROM golang:1.22 AS builder

WORKDIR /app

# Копируем go.mod сервиса
COPY go.mod go.sum ./
RUN go mod download

# Копируем ВСЮ папку сервиса (internal, gen, cmd)
COPY . .

WORKDIR /app/cmd/auth-service

RUN go build -o auth-service .

# ---------- RUNTIME ----------
FROM alpine:3.18

WORKDIR /app

COPY --from=builder /auth-service .

EXPOSE 50051

CMD ["./auth-service"]
