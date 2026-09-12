.PHONY: run build test migrate tidy

run:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

test:
	go test ./...

tidy:
	go mod tidymigrate:
	@echo "apply migrations/0001_init.sql against DATABASE_URL=$$DATABASE_URL"