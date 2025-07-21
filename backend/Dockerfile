# Многоэтапная сборка для оптимизации размера образа
FROM golang:1.23-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарный файл
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Финальный образ
FROM alpine:latest

# Устанавливаем ca-certificates для HTTPS запросов
RUN apk --no-cache add ca-certificates tzdata

# Создаем непривилегированного пользователя
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /root/

# Копируем бинарный файл из builder stage
COPY --from=builder /app/main .

# Копируем миграции
COPY --from=builder /app/migrations ./migrations

# Меняем владельца файлов
RUN chown -R appuser:appgroup /root/

# Переключаемся на непривилегированного пользователя
USER appuser

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]