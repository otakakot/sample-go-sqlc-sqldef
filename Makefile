SHELL := /bin/bash
include .env
export
export APP_NAME := $(basename $(notdir $(shell pwd)))

.PHONY: help
help: ## display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: up
up: ## docker compose up with air hot reload
	@docker compose --project-name ${APP_NAME} --file ./.docker/compose.yaml up -d

.PHONY: down
down: ## docker compose down
	@docker compose --project-name ${APP_NAME} down --volumes

.PHONY: balus
balus: ## destroy everything about docker. (containers, images, volumes, networks.)
	@docker compose --project-name ${APP_NAME} down --rmi all --volumes

.PHONY: psql
psql:
	@docker exec -it postgres psql -U postgres

.PHONY: mig
mig: ## run migration
	@psqldef --user=postgres --password=postgres --host=localhost --port=5432 postgres < schema/schema.sql 

.PHONY: gen
gen: ## generate sqlboiler
	@find pkg/schema -type f -not -name "*.sql" -exec rm -rf {} \;
	@sqlc generate
	@go mod tidy
	@go mod vendor
