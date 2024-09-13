build:
	go build -o bin/${BINARY} ./cmd/api/main.go  

run:
	go run cmd/api/main.go

scratch:
	go run cmd/scratch/scratch.go

up:
	@echo "Starting docker containers..."
	docker compose up --build -d --remove-orphans

down:
	@echo "Stopping containers..."
	docker compose down