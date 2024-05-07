.PHONY: build-all
build-all:
	cd cart && make build
	cd loms && make build
	cd notifier && make build

.PHONY: run-all
run-all:
	docker-compose up --force-recreate --build -V

.PHONY: test-all
test-all:
	cd cart && make test
	cd loms && make test
	cd notifier && make test

.PHONY: lint-all
lint-all:
	cd cart && make lint
	cd loms && make lint
	cd notifier && make lint

.PHONY: e2e
e2e:
	go test ./e2etest/cart -count=1
