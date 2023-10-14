include .env

LOCAL_BIN:=$(CURDIR)/bin
LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(POSTGRES_PORT) dbname=$(POSTGRES_DB) user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) sslmode=disable"

migration-down:
	goose --dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

migration-reset:
	goose --dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} reset -v

install-golangci-lint:
	GOBIN=${LOCAL_BIN} go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54

install-deps:
	GOBIN=${LOCAL_BIN} go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=${LOCAL_BIN} go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

generate:
	make generate-note-api

generate-user-api:
	mkdir -p pkg/chat_v1
	protoc --proto_path api/chat_v1 \
	--go_out=pkg/chat_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/chat_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/chat_v1/chat.proto

lint:
	./bin/golangci-lint run ./... --config .golangci.pipeline.yaml

build:
	GOOS=linux GOARCH=amd64 go build -o ./bin/grpc-server ./cmd/grpc-server/main.go

docker-build-and-push:
	docker buildx build --no-cache --platform linux/amd64 -t cr.selcloud.ru/ontropos42/test-server:v0.0.1 .
	docker login -u token -p CRgAAAAAT7M2IVc1bUBai6HdzbxITsRZGKhct7XO cr.selcloud.ru/ontropos42
	docker push cr.selcloud.ru/ontropos42/test-server:v0.0.1