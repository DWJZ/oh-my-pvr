.PHONY: help build-server build-client build compose-up compose-down

help:
	@echo "Available commands:"
	@echo "   make build-server    - Build the Go server Docker image"
	@echo "   make build-client    - Build the Next.js client Docker image (using pnpm)"
	@echo "   make build           - Build both server and client images"
	@echo "   make compose-up      - Build and start both containers with Docker Compose"
	@echo "   make compose-down    - Stop Docker Compose containers"

# Build the Go backend Docker image.
build-server:
	docker build -t omp-server ./server

# Build the Next.js client Docker image using pnpm.
build-client:
	docker build -t omp-client ./client

# Build both server and client images.
build: docker compose build

# Run both containers using Docker Compose.
compose-up:
	docker compose up -d

# Stop and remove containers started via Docker Compose.
compose-down:
	docker compose down
