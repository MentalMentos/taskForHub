GATEWAY_NAME=api-gateway
AUTH_NAME=auth
BOOKS_NAME=books
MONGO_NAME=mongo

MONGO_IMAGE=mongo:latest

run: build
	docker-compose up --build

build:
	docker-compose build

# Остановка контейнеров
stop:
	docker-compose down

# Очистка всех остановленных контейнеров
clean:
	docker-compose down --volumes --remove-orphans

gateway:
	docker-compose up -d $(GATEWAY_NAME)

auth:
	docker-compose up -d $(AUTH_NAME)

books:
	docker-compose up -d $(BOOKS_NAME)

# Запуск контейнера только для MongoDB
mongo:
	docker-compose up -d $(MONGO_NAME)

restart:
	docker-compose down
	docker-compose up --build
