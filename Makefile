.PHONY: build test

build:
	go build -o bin/edged ./cmd/edged

test:
	go test ./...
