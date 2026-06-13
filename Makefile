.PHONY: build test tidy

build:
	go build ./cmd/cognitor

test:
	go test ./...

tidy:
	go mod tidy
