# Builder stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./api-gateway ./api-gateway

RUN go build -o /app/api-gateway ./api-gateway

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/api-gateway .

CMD ["./api-gateway"]
