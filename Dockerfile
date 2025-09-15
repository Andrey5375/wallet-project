FROM golang:1.24-alpine

WORKDIR /app

# Установка psql и bash
RUN apk add --no-cache postgresql-client bash

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN chmod +x wait-for-postgres.sh
RUN go build -o server ./cmd/server

EXPOSE 8080

CMD ["./server"]
