include .env
export

export PROJECT_ROOT=$(shell pwd)

local-run:
	@export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run cmd/tg-bot/main.go

# PostgreSQL
postgres-up:
	@docker compose up -d postgres

postgres-down:
	@docker compose down postgres

# Migrations
migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Required parameter seq is missing. \nExample: make migrate-create seq=init"; \
		exit 1; \
	fi; \

	@docker compose run --rm postgres-migrate \
	 create \
	 -ext sql \
	 -dir /migrations \
	 -seq "$(seq)"

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Required parameter 'action' is missing. \nAvailable actions:\nup (count)\ndown (count)\nversion\nforce (count)"; \
		exit 1; \
	fi; \

	@docker compose run --rm postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable \
		$(action)

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

# Logs
logs-cleanup:
	@read -p "WARN: Do you want to delete all log files files? [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		rm -rf ${PROJECT_ROOT}/logs && \
		echo "Log files cleared"; \
	else \
		echo "Cancelled"; \
	fi