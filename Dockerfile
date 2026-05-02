# Этап 1: Сборка
FROM golang:1.23-alpine AS builder

# Заставляем Go использовать локальную версию 1.23
ENV GOTOOLCHAIN=local

WORKDIR /app
COPY go.mod go.sum ./

# На всякий случай фиксируем версию в go.mod прямо при сборке
RUN sed -i 's/go 1.25/go 1.23/g' go.mod

RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Этап 2: Запуск
FROM alpine:latest
WORKDIR /app

# Устанавливаем пакет для работы с таймзонами (Asia/Almaty)
RUN apk add --no-cache tzdata

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080
CMD ["./main"]