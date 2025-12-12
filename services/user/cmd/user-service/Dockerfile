FROM golang:1.22 AS builder

WORKDIR /app

# Копируем всю монорепу (go.work использует все модули)
COPY ../.. .

# Переходим в папку сервиса
WORKDIR /app/services/user/cmd/user-service

RUN go mod tidy

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o /user-service .

# ——— Runtime image ———
FROM alpine:3.18

WORKDIR /root/

# Копируем бинарник из builder-слоя
COPY --from=builder /user-service .

# Стартуем сервис
CMD ["./user-service"]
