LOCAL_BIN:=$(CURDIR)/bin
GOLANGCI_BIN:=$(LOCAL_BIN)/golangci-lint
PROTOC_BIN:=$(LOCAL_BIN)/protoc

PHONY: build
build:
	@go build -o bin/cart-server ./cmd/cart-server

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
	go test ./... -cover -race -parallel 2

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

# proto для сервиса loms
vendor/loms:
	mkdir -p vendor/loms
	cp  ../loms/api/v1/loms.proto vendor/loms/loms.proto

# Вендоринг внешних proto файлов
.vendor: vendor/google/protobuf vendor/validate vendor/google/api vendor/protoc-gen-openapiv2/options vendor/loms

.PHONY: .bin-deps
.bin-deps:
	$(info Installing binary dependencies...)
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1 && \
    GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

LOMS_PROTO_PATH=vendor/loms/loms.proto
LOMS_OUT_PATH=internal/pkg/pb

PHONY: protoc-generate
protoc-generate: .bin-deps .vendor
	mkdir -p $(LOMS_OUT_PATH)
	protoc \
	-I vendor \
	--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
	--go_out $(LOMS_OUT_PATH) \
	--go_opt paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
	--go-grpc_out $(LOMS_OUT_PATH) \
	--go-grpc_opt paths=source_relative \
	$(LOMS_PROTO_PATH)
	go mod tidy
