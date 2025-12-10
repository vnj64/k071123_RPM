ifneq (,$(wildcard .env))
    include .env
    export
endif

SERVICE ?= parking

TEST_PACKAGES := $(shell go list ./internal/services/$(SERVICE)_service/... | grep -v '/mocks')

COVERAGE_FILE := coverage.out
HTML_COVERAGE_FILE := coverage.html

.PHONY: test
test:
	@echo "Testing service: $(SERVICE)"
	go test ./internal/services/$(SERVICE)_service/... \
    -coverprofile=coverage.out \
    -coverpkg=./internal/services/$(SERVICE)_service/...

.PHONY: cover
cover: test
	@echo "Generating HTML coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage HTML saved to $(HTML_COVERAGE_FILE)"

swag:
	go run tools/swaggergen.go

migrate-parking-up:
	goose -dir internal/services/parking_service/storage/migrations postgres "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(PARKING_DB_NAME)" up

migrate-notification-up:
	goose -dir ./internal/services/notification_service/storage/migrations postgres "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(NOTIFICATION_DB_NAME)" up

migrate-user-up:
	goose -dir ./internal/services/user_service/storage/migrations postgres "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(USER_DB_NAME)" up

migrate-order-up:
	goose -dir ./internal/services/order_service/storage/migrations postgres "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(ORDER_DB_NAME)" up

migrate-all-up: migrate-parking-up migrate-user-up migrate-order-up

migrate-parking-down:
	goose -dir ./internal/services/parking_service/storage/migrations postgres "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(PARKING_DB_NAME)" down

migrate-notification-down:
	goose -dir ./internal/services/notification_service/storage/migrations postgres "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(NOTIFICATION_DB_NAME)" down

migrate-user-down:
	goose -dir ./internal/services/user_service/storage/migrations postgres postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(USER_DB_NAME) down

migrate-order-down:
	goose -dir ./internal/services/order_service/storage/migrations postgres "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(ORDER_DB_NAME)" down


migrate-all-down: migrate-parking-down migrate-notification-down migrate-user-down

# Other
# Реализовывать в корне notification-service
gen-notify-grpc:
	protoc -I=proto \
      --go_out=. \
      --go-grpc_out=. \
      proto/notification_service.proto
