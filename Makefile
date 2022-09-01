.PHONY: build-default
build-default:
	go build -o ./bin/home-library -v ./cmd/home-library

.PHONY: run-local
run-local: build-default
	./bin/home-library --env=local

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := build-default
