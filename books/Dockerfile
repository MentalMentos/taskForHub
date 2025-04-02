# Builder stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./books ./books

RUN go build -o /books ./books/cmd

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /books .

CMD ["./books"]
