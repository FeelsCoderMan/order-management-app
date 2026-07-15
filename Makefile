BINARY_DIR := ./dist
CMD_PATH   := ./cmd
BINARY     := $(BINARY_DIR)/server
GO_FLAGS   := -ldflags="-s -w"

.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make run       Run the application locally"
	@echo "  make build     Build the application"
	@echo "  make clean     Remove build artifacts"
	@echo "Docker:"
	@echo "  make up        Build and start all containers"
	@echo "  make start     Start existing containers"
	@echo "  make stop      Stop containers"
	@echo "  make down      Stop and remove containers"
	@echo "  make down-v    Stop, remove containers and volumes"
	@echo "  make restart   Restart containers"
	@echo "  make rebuild   Rebuild application image and start containers"
	@echo "  make logs      View container logs"
	@echo "  make ps        List running containers"

run:
	go run $(CMD_PATH)

build:
	go build $(GO_FLAGS) -o $(BINARY) $(CMD_PATH)

clean:
	rm -rf $(BINARY_DIR)

up:
	docker compose up --build -d

start:
	docker compose start

stop:
	docker compose stop

down:
	docker compose down

down-v:
	docker compose down -v

restart:
	docker compose restart

rebuild: down up

logs:
	docker compose logs -f

ps:
	docker compose ps
