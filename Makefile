include db.env

export MIGRATE_PATH = ./db/migrations

BUILD_PATH := ./build
DOCKER_COMPOSE_PATH := $(BUILD_PATH)/docker-compose.yml
DOCKER_COMPOSE_LOCAL_PATH := $(BUILD_PATH)/docker-compose.local.yml

# Docker Compose
build:
	@docker-compose -f $(DOCKER_COMPOSE_PATH) build

up:
	@docker-compose -f $(DOCKER_COMPOSE_PATH) up -d

down:
	@docker-compose -f $(DOCKER_COMPOSE_PATH) down

logs:
	@docker-compose -f $(DOCKER_COMPOSE_PATH) logs

local:
	@docker-compose -f $(DOCKER_COMPOSE_LOCAL_PATH) up -d

local-down:
	@docker-compose -f $(DOCKER_COMPOSE_LOCAL_PATH) down

local-logs:
	@docker-compose -f $(DOCKER_COMPOSE_LOCAL_PATH) logs

# Migrations
migrate-create:
	@bash ./scripts/migrate_create.sh

migrate-up:
	@migrate \
		-database $(DB_URL) \
		-path $(MIGRATE_PATH) \
		up

migrate-down:
	@migrate \
		-database $(DB_URL) \
		-path $(MIGRATE_PATH) \
		down

# Swagger
swagger-generate:
	@swag init -g ./cmd/app/main.go

.PHONY: build up down logs local local-down local-logs \
		migrate-create migrate-up migrate-down \
		swagger-generate
