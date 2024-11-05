include local.env

LOCAL_BIN:=$(CURDIR)/bin


#LOCAL_MIGRATION_DSN=${MIGRATION_DSN_MAKE}

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v1.0.4


get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc


generate:
	make generate-chat-api

generate-chat-api:
	mkdir -p pkg/chat_v1
	protoc --proto_path=api/chat_v1 --proto_path=vendor.protogen \
	--go_out=pkg/chat_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/chat_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	--validate_out=lang=go:pkg/chat_v1 --validate_opt=paths=source_relative \
	--plugin=protoc-gen-validate=bin/protoc-gen-validate \
	api/chat_v1/chat.proto

vendor-proto:
	@if [ ! -d vendor.protogen/validate ]; then \
		mkdir -p vendor.protogen/validate &&\
		git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/protoc-gen-validate &&\
		mv vendor.protogen/protoc-gen-validate/validate/*.proto vendor.protogen/validate &&\
		rm -rf vendor.protogen/protoc-gen-validate ;\
	fi

# Для локальной накатки миграций
# local-migration-status:
# 	bin/goose -dir ${MIGRATION_DIR}	postgres ${LOCAL_MIGRATION_DSN} status -v

# local-migration-up:
# 	bin/goose -dir ${MIGRATION_DIR}	postgres ${LOCAL_MIGRATION_DSN} up -v

# local-migration-down:
# 	bin/goose -dir ${MIGRATION_DIR}	postgres ${LOCAL_MIGRATION_DSN} down -v

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml

test: 
	go clean -testcache
	go test ./... -covermode count -coverpkg=github.com/Mobo140/microservices/chat/internal/service/..., github.com/Mobo140/microservices/chat/internal/transport/... -count 5

test-coverage:
	go clean -testcache
	go test ./... -coverprofile=coverage.tmp.out -covermode count -coverpkg=github.com/Mobo140/microservices/chat/internal/transport/handlers/chat/...,github.com/Mobo140/microservices/chat/internal/service/chat/... -count 5
	grep -v 'mocks\|config' coverage.tmp.out > coverage.out 
	rm coverage.tmp.out
	go tool cover -html=coverage.out
	go tool cover -func=./coverage.out | grep "total"
	grep -sqFx "/coverage.out" .gitignore || echo "coverage_out" >> .gitignore