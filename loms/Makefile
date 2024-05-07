LOCAL_BIN:=$(CURDIR)/bin
GOLANGCI_BIN:=$(LOCAL_BIN)/golangci-lint
PROTOC_BIN:=$(LOCAL_BIN)/protoc
MIGRATIONS_DIR = ./migrations
PG_MASTER_DSN = postgresql://postgres:myozon@localhost:5432/postgres

PHONY: build
build:
	@go build -o bin/loms-server ./cmd/loms-server

.PHONY: .install-lint
.install-lint:
ifeq ($(wildcard $(GOLANGCI_BIN)),)
	$(info Downloading golangci-lint)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LOCAL_BIN) v1.56.2
endif

.PHONY: lint
lint: .install-lint
	$(info Running lint...)
	$(GOLANGCI_BIN) run -v --fix

.PHONY: test
test:
	go test ./... -cover

.PHONY: .install-goose
.install-goose:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest

.PHONY: migrate
migrate: .install-goose
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(PG_MASTER_DSN) GOOSE_MIGRATION_DIR=$(MIGRATIONS_DIR) $(LOCAL_BIN)/goose up

.PHONY: .install-sqlc
.install-sqlc:
	GOBIN=$(LOCAL_BIN) go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

.PHONY: sqlc-generate
sqlc-generate: .install-sqlc
	$(LOCAL_BIN)/sqlc generate

# proto описания google/protobuf
vendor/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor/protobuf &&\
	cd vendor/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	mkdir -p vendor/google
	mv vendor/protobuf/src/google/protobuf vendor/google
	rm -rf vendor/protobuf


# proto описания validate
vendor/validate:
	git clone -b main --single-branch --depth=2 --filter=tree:0 \
		https://github.com/bufbuild/protoc-gen-validate vendor/tmp && \
		cd vendor/tmp && \
		git sparse-checkout set --no-cone validate &&\
		git checkout
		mkdir -p vendor/validate
		mv vendor/tmp/validate vendor/
		rm -rf vendor/tmp


# proto описания google/googleapis
vendor/google/api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis vendor/googleapis && \
 	cd vendor/googleapis && \
	git sparse-checkout set --no-cone google/api && \
	git checkout
	mkdir -p  vendor/google
	mv vendor/googleapis/google/api vendor/google
	rm -rf vendor/googleapis

# proto описания protoc-gen-openapiv2/options
vendor/protoc-gen-openapiv2/options:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway vendor/grpc-ecosystem && \
 	cd vendor/grpc-ecosystem && \
	git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
	git checkout
	mkdir -p vendor/protoc-gen-openapiv2
	mv vendor/grpc-ecosystem/protoc-gen-openapiv2/options vendor/protoc-gen-openapiv2
	rm -rf vendor/grpc-ecosystem

# Вендоринг внешних proto файлов
.vendor: vendor/google/protobuf vendor/validate vendor/google/api vendor/protoc-gen-openapiv2/options

.PHONY: .bin-deps
.bin-deps:
	$(info Installing binary dependencies...)
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1 && \
    GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0 && \
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v1.0.4 && \
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.19.1 && \
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.19.1 && \

LOMS_PROTO_PATH=api/v1/loms.proto
LOMS_OUT_PATH=internal/pkg/pb/loms

PHONY: protoc-generate
protoc-generate: .bin-deps .vendor
	mkdir -p $(LOMS_OUT_PATH)
	protoc \
	-I api/v1 \
	-I vendor \
	--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
	--go_out $(LOMS_OUT_PATH) \
	--go_opt paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
	--go-grpc_out $(LOMS_OUT_PATH) \
	--go-grpc_opt paths=source_relative \
	--plugin=protoc-gen-validate=$(LOCAL_BIN)/protoc-gen-validate \
	--validate_out="lang=go,paths=source_relative:$(LOMS_OUT_PATH)" \
	--plugin=protoc-gen-grpc-gateway=$(LOCAL_BIN)/protoc-gen-grpc-gateway \
	--grpc-gateway_out $(LOMS_OUT_PATH) \
	--grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true \
	--plugin=protoc-gen-openapiv2=$(LOCAL_BIN)/protoc-gen-openapiv2 \
	--openapiv2_out $(LOMS_OUT_PATH) \
	--openapiv2_opt logtostderr=true \
	$(LOMS_PROTO_PATH)
	go mod tidy
