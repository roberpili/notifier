.PHONY: help up down build logs restart clean test lint vulncheck swagger
include .env
export

help:
	@echo "Available commands:"
	@echo "  make up       		- start all services"
	@echo "  make down     		- stop all services"
	@echo "  make build    		- build all images"
	@echo "  make logs     		- view all logs"
	@echo "  make restart  		- restart all services"
	@echo "  make clean    		- stop and remove everything"
	@echo "  make test     		- run backend tests"
	@echo "  make lint     		- run linter on backend code"
	@echo "  make vulncheck     - run vulnerabilities check on backend code"
	@echo "  make swagger     	- generate openapi documentation"

up:
	docker compose up -d

down:
	docker compose down --remove-orphans

build:
	docker compose build

logs:
	docker compose logs -f

restart: down build up

clean:
	docker compose down -v --rmi all

test:
	cd notifier-backend && go test ./... -v

lint:
	golangci-lint run --timeout 10m --config staticcheck.yml


vulncheck:
	govulncheck ./...

swagger:
	swag init -g cmd/main.go -o docs

MIGRATE=migrate -path ./migrations -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable"

migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

migrate-force:
	$(MIGRATE) force 3
