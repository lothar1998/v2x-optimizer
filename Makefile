CGO_ENABLED = 0

MOCKS_DIR = ./test/mocks
MOCK_DEFINITION_FILE = mocks.go

.PHONY: all build lint test format install mocks optimize-cplex
all: mocks lint test build

# V2X Optimizer
build:
	@echo "# Running build..."
	@go build ./...

lint:
	@echo "# Running lint..."
	@golangci-lint run

test:
	@echo "# Running test..."
	@go test ./...

format:
	@echo "# Running format..."
	@goimports -w pkg/
	@goimports -w internal/
	@goimports -w cmd/

mocks:
	@echo "# Regenerating mocks..."
	@find $(MOCKS_DIR) -type f | grep -v $(MOCK_DEFINITION_FILE) | xargs rm
	@go generate $(MOCKS_DIR)/$(MOCK_DEFINITION_FILE)

install:
	@echo "# Installing app..."
	@go install ./...

# CPLEX optimizer
optimize-cplex:
	@echo "# Optimizing using CPLEX..."
	@$(MAKE) --no-print-directory -C third_party/cplex
