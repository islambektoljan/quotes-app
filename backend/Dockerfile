FROM golang:1.25-alpine

WORKDIR /app

# Копируем go.mod и ВСЕ исходники
COPY go.mod ./
COPY . .

# Генерируем go.sum
RUN go mod tidy

# Скачиваем зависимости
RUN go mod download

# Собираем приложение
RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
