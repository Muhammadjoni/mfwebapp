DB_DSN ?= postgres://mfuser:mfpass@localhost:5432/mfwebapp?sslmode=disable

.PHONY: db db-stop backend frontend-user frontend-admin dev build

## Start PostgreSQL and Redis in Docker
db:
	docker-compose up postgres redis -d

## Stop all containers
db-stop:
	docker-compose down

## Run backend in dev mode (migrations run automatically on startup)
backend:
	cd backend && go run ./cmd/api

## Build backend binary
build:
	cd backend && go build -o bin/api ./cmd/api

## Run frontend-user dev server
frontend-user:
	cd frontend-user && npm run dev

## Run frontend-admin dev server
frontend-admin:
	cd frontend-admin && npm run dev

## Install frontend deps
install:
	cd frontend-user && npm install
	cd frontend-admin && npm install

## Run all services via docker-compose
up:
	docker-compose up --build

## go vet + build check
lint:
	cd backend && go vet ./...
