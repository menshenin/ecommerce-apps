LOCAL_BIN:=$(CURDIR)/bin
GOLANGCI_BIN:=$(LOCAL_BIN)/golangci-lint

PHONY: build
build:
	@go build -o bin/consumer ./cmd/consumer

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
