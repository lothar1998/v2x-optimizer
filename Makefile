CGO_ENABLED = 0
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOHOSTARCH)
SRC_DIRS = ./cmd/... ./pkg/...

.PHONY: all build lint test format
all: info build lint test

info:
	@pwd
	@echo $(GOPATH)
	@echo $(GOBIN)
	go env

build:
	@echo "Running build"
	go build -mod vendor ./...

lint:
	@echo "Running lint"
	golint -set_exit_status $(SRC_DIRS)

test:
	@echo "Running test"
	go test -mod vendor $(SRC_DIRS)

format:
	@echo "Running format"
	go fmt $(SRC_DIRS)
