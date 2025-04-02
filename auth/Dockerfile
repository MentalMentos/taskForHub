# Builder stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./auth ./auth

RUN go build -o /auth ./auth/cmd


# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /auth .

CMD ["./auth"]
