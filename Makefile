VERSION  ?= dev
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)

.PHONY: all build test test-cover lint vet fmt docker-up docker-down clean

all: lint test build

## Build
build:
	CGO_ENABLED=0 go build -ldflags "-s -w \
		-X main.Version=$(VERSION) \
		-X main.BuildTime=$$(date -u +%Y-%m-%dT%H:%M:%SZ) \
		-X main.GitCommit=$(GIT_COMMIT)" \
		-o bin/api ./cmd/api

## Test
test:
	go test ./... -race -count=1

test-cover:
	go test ./... -race -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

## Code Quality
lint:
	golangci-lint run ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

## Docker
docker-up:
	docker compose up --build -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

## Misc
clean:
	rm -rf bin/ coverage.out coverage.html
