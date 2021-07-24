CGO_ENABLED = 0
SRC_DIRS = ./cmd/... ./pkg/...

.PHONY: all build lint test format
all: build lint test

build:
	@echo "Running build"
	go build ./...

lint:
	@echo "Running lint"
	golint -set_exit_status ./...

test:
	@echo "Running test"
	go test ./...

format:
	@echo "Running format"
	go fmt ./...
