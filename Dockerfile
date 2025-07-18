FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO отключен для совместимости с alpine
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/app .

RUN chmod +x ./app

EXPOSE 8080

CMD ["./app"]
