.DEFAULT_GOAL := help

# Setup the application
setup:
	@chmod +x setup.sh
	@./setup.sh

# Run the application locally
run: setup
	@go run ./cmd/api

# Build the binary locally
build:
	@echo "Building binary..."
	@go build -o main ./cmd/api

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean and tidy go modules
tidy:
	@echo "Tidying go modules..."
	@go mod tidy

# Clean built files
clean:
	@echo "Cleaning up..."
	@rm -f main

# Live Reload with Air
watch: setup
	@if command -v air > /dev/null; then \
		air; \
	else \
		read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/air-verse/air@latest; \
			air; \
		else \
			echo "You chose not to install air. Exiting..."; \
			exit 1; \
		fi; \
	fi

# Docker: Build containers
docker-build: setup
	@docker compose build

# Docker: Build and start all containers in background
docker-run: setup
	@docker compose up --build -d

# Docker: Stop all containers
docker-down:
	@docker compose down

# Docker: Stop all containers and remove volumes (clean database reset)
docker-down-v:
	@docker compose down -v

# Docker: Follow application logs
docker-logs:
	@docker compose logs -f carrpigeon

# Docker: Restart application container
docker-restart:
	@docker compose restart carrpigeon

# Help
help:
	@echo "Available commands:"
	@echo "  make setup          - Initialize configs and storage"
	@echo "  make run            - Run application locally"
	@echo "  make build          - Build binary locally"
	@echo "  make watch          - Run with live-reload (air)"
	@echo "  make test           - Run tests"
	@echo "  make tidy           - Run go mod tidy"
	@echo "  make clean          - Remove built binary"
	@echo "  make docker-build   - Build Docker images"
	@echo "  make docker-run     - Build and start all containers in background"
	@echo "  make docker-down    - Stop all containers"
	@echo "  make docker-down-v  - Stop containers and remove volumes (database reset)"
	@echo "  make docker-logs    - Follow application logs in Docker"
	@echo "  make docker-restart - Restart application container"

.PHONY: setup build run test tidy clean watch docker-build docker-run docker-down docker-down-v docker-logs docker-restart help
