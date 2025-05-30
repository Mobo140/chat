include env/local.env

LOCAL_BIN := $(CURDIR)/bin

# Setup and run project 
setup: install-deps generate up
	
run: 
	go run cmd/grpc-server/main.go --config-path=env/local.env -l=debug

# Start all services in detached mode
up:
	docker-compose up -d

down:
	docker-compose down

# Install CLI tools needed for protobuf, migrations, grpc-gateway, lint
install-deps:
	@test -f $(LOCAL_BIN)/protoc-gen-go || GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	@test -f $(LOCAL_BIN)/protoc-gen-go-grpc || GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	@test -f $(LOCAL_BIN)/goose || GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0
	@test -f $(LOCAL_BIN)/protoc-gen-validate || GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v1.0.4
	@test -f $(LOCAL_BIN)/protoc-gen-grpc-gateway || GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.20.0
	@test -f $(LOCAL_BIN)/protoc-gen-openapiv2 || GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.20.0
	@test -f $(LOCAL_BIN)/statik || GOBIN=$(LOCAL_BIN) go install github.com/rakyll/statik@v0.1.7
	@test -f $(LOCAL_BIN)/golangci-lint || GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0

# Generate protobuf, grpc, validation, gateway and swagger code
generate:
	mkdir -p pkg/swagger pkg/chat_v1
	make generate-chat-api
	$(LOCAL_BIN)/statik -src=pkg/swagger/ -include='*.css,*.html,*.js,*.json,*.png'

generate-chat-api:
	protoc --proto_path=api/chat_v1 --proto_path=vendor.protogen \
		--go_out=pkg/chat_v1 --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=bin/protoc-gen-go \
		--go-grpc_out=pkg/chat_v1 --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
		--validate_out=lang=go:pkg/chat_v1 --validate_opt=paths=source_relative \
		--plugin=protoc-gen-validate=bin/protoc-gen-validate \
		--grpc-gateway_out=pkg/chat_v1 --grpc-gateway_opt=paths=source_relative \
		--plugin=protoc-gen-grpc-gateway=bin/protoc-gen-grpc-gateway \
		--openapiv2_out=allow_merge=true,merge_file_name=api:pkg/swagger \
		--plugin=protoc-gen-openapiv2=bin/protoc-gen-openapiv2 \
		api/chat_v1/chat.proto

# Clone third-party proto deps if missing
vendor-proto:
	@if [ ! -d vendor.protogen/validate ]; then \
		mkdir -p vendor.protogen/validate && \
		git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/protoc-gen-validate && \
		mv vendor.protogen/protoc-gen-validate/validate/*.proto vendor.protogen/validate && \
		rm -rf vendor.protogen/protoc-gen-validate ; \
	fi
	@if [ ! -d vendor.protogen/google ]; then \
		git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis && \
		mkdir -p vendor.protogen/google && \
		mv vendor.protogen/googleapis/google/api vendor.protogen/google && \
		rm -rf vendor.protogen/googleapis ; \
	fi
	@if [ ! -d vendor.protogen/protoc-gen-openapiv2 ]; then \
		mkdir -p vendor.protogen/protoc-gen-openapiv2/options && \
		git clone https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/openapiv2 && \
		mv vendor.protogen/openapiv2/protoc-gen-openapiv2/options/*.proto vendor.protogen/protoc-gen-openapiv2/options && \
		rm -rf vendor.protogen/openapiv2 ; \
	fi

# Generate self-signed certs for local TLS testing
gen-cert:
	mkdir -p secure
	openssl genrsa -out secure/ca.key 4096
	openssl req -new -x509 -key secure/ca.key -sha256 -subj "/C=US/ST=NJ/O=CA, Inc." -days 365 -out secure/ca.cert
	openssl genrsa -out secure/service.key 4096
	openssl req -new -key secure/service.key -out secure/service.csr -config certificate.conf
	openssl x509 -req -in secure/service.csr -CA secure/ca.cert -CAkey secure/ca.key -CAcreateserial \
		-out secure/service.pem -days 365 -sha256 -extfile certificate.conf -extensions req_ext

# Run linter with pipeline config
lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml

# Run tests with coverage and multiple retries
test:
	go clean -testcache
	go test ./... -covermode=count -coverpkg=github.com/Mobo140/microservices/chat/internal/service/...,github.com/Mobo140/microservices/chat/internal/transport/... -count=5

# Run tests with coverage report generation and open HTML
test-coverage:
	go clean -testcache
	go test ./... -coverprofile=coverage.tmp.out -covermode=count -coverpkg=github.com/Mobo140/microservices/chat/internal/transport/handlers/chat/...,github.com/Mobo140/microservices/chat/internal/service/chat/... -count=5
	grep -v 'mocks\|config' coverage.tmp.out > coverage.out
	rm coverage.tmp.out
	go tool cover -html=coverage.out
	go tool cover -func=coverage.out | grep "total"
	grep -sqFx "/coverage.out" .gitignore || echo "coverage_out" >> .gitignore

# Run gRPC load test with TLS certs
grpc-load-test:
	ghz --proto api/test_chat_v1/test_chat.proto \
		--call chat_v1.ChatV1.Get \
		--data '{"id": 1}' \
		--rps 100 --total 3000 \
		--cacert ca.cert --cert service.pem --key service.key \
		localhost:${GRPC_PORT}

