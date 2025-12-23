FROM golang:1.22 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o search-service ./cmd/search-service/

FROM alpine:3.18
WORKDIR /root/
COPY --from=builder /app/search-service .

EXPOSE 8080
EXPOSE 50058
CMD ["./search-service"]
