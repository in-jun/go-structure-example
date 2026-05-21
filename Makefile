SERVICES  := auction bid payment gateway auth
VERSION   ?= dev
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)

.PHONY: all build test test-cover lint vet fmt docker-build docker-up docker-down clean proto migrate-up

all: lint test build

## Build
build:
	@for svc in $(SERVICES); do \
		echo "building $$svc..."; \
		CGO_ENABLED=0 go build -ldflags "-s -w \
			-X main.Version=$(VERSION) \
			-X main.BuildTime=$$(date -u +%Y-%m-%dT%H:%M:%SZ) \
			-X main.GitCommit=$(GIT_COMMIT)" \
			-o bin/$$svc ./cmd/$$svc; \
	done

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
docker-build:
	docker compose build --build-arg VERSION=$(VERSION) --build-arg GIT_COMMIT=$(GIT_COMMIT)

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

## Database
migrate-up:
	@for svc in auction bid payment auth; do \
		echo "migrating $$svc..."; \
		migrate -path migrations/$$svc -database "postgres://postgres:postgres@localhost:5432/$${svc}_db?sslmode=disable" up; \
	done

## Proto
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/auction/v1/auction.proto

## Misc
health:
	@curl -s localhost:8080/health/live | jq .
	@curl -s localhost:8080/health/ready | jq .

clean:
	rm -rf bin/ coverage.out coverage.html
