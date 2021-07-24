CGO_ENABLED = 0
SRC_DIRS = ./cmd/... ./pkg/...

.PHONY: all build lint test format
all: build lint test

build:
	@echo "Running build"
	go build ./...

lint:
	@echo "Running lint"
	go list ./... | grep -v /vendor/ | xargs -L1 golint -set_exit_status

test:
	@echo "Running test"
	go test $(SRC_DIRS)

format:
	@echo "Running format"
	go fmt $(SRC_DIRS)
