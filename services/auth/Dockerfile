# ---------- BUILD ----------
FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service ./cmd/auth-service/

# ---------- RUN ----------
FROM alpine:3.18

WORKDIR /root/

COPY --from=builder /app/auth-service .

EXPOSE 8080
EXPOSE 50051

CMD ["./auth-service"]
