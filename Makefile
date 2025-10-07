SHELL := bash

.EXPORT_ALL_VARIABLES:

-include .env

# Define the default target
.PHONY: default
default: help

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  migrate-up - Process to up migrating database"
	@echo "  migrate-up - Process to rollback migrating database"
	@echo "  seeder FILENAME - Process the specified seeder file name"

run-dev:
	echo "Starting Application In Development Mode"
	go run ./main.go

install:
	go mod download

migrate-up: 
	migrate -database "postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -path database/migrations up


migrate-down: 
	migrate -database "postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -path database/migrations down