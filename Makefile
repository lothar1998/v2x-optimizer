CGO_ENABLED = 0

.PHONY: all build lint test format install
all: lint test build

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

install:
	@echo "Installing app"
	go install ./...
