ifneq (,$(wildcard .env))
    include .env
    export
endif

swag-user:
	swag init --parseDependency -g internal/services/user_service/cmd/main.go --output=./internal/services/user_service/docs

swag-notify:
	swag init --parseDependency -g internal/services/notification_service/cmd/main.go --output=./internal/services/notification_service/docs

migrate-up:
	goose up

migrate-down:
	goose down

# Реализовывать в корне notification-service
gen-notify-grpc:
	protoc -I=proto \
      --go_out=. \
      --go-grpc_out=. \
      proto/notification_service.proto
