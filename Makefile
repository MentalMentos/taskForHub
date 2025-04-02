include .env

up:
	docker-compose up --build -d
	go run cmd/app/main.go
down:
	docker-compose down

run: up