include local.env

LOCAL_BIN:=$(CURDIR)/bin


#LOCAL_MIGRATION_DSN=${MIGRATION_DSN_MAKE}

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc


generate:
	make generate-chat-api

generate-chat-api:
	mkdir -p pkg/chat_v1
	protoc --proto_path api/chat_v1 \
	--go_out=pkg/chat_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/chat_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/chat_v1/chat.proto

# Для локальной накатки миграций
# local-migration-status:
# 	bin/goose -dir ${MIGRATION_DIR}	postgres ${LOCAL_MIGRATION_DSN} status -v

# local-migration-up:
# 	bin/goose -dir ${MIGRATION_DIR}	postgres ${LOCAL_MIGRATION_DSN} up -v

# local-migration-down:
# 	bin/goose -dir ${MIGRATION_DIR}	postgres ${LOCAL_MIGRATION_DSN} down -v

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml