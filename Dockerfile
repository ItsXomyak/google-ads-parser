# syntax=docker/dockerfile:1.4
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Установим зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Собираем бинарник
RUN go build -o app ./cmd

# Минимальный рантайм-образ
FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]
