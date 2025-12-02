ifneq (,$(wildcard .env))
    include .env
    export
endif

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
