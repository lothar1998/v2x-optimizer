CGO_ENABLED = 0

.PHONY: all build lint test format install optimize
all: lint test build

# V2X Optimizer
build:
	@echo "# Running build..."
	@go build ./...

lint:
	@echo "# Running lint..."
	@golint -set_exit_status ./...

test:
	@echo "# Running test..."
	@go test ./...

format:
	@echo "# Running format..."
	@go fmt ./...

install:
	@echo "# Installing app..."
	@go install ./...

# CPLEX optimizer
optimize-cplex:
	@echo "# Optimizing using CPLEX..."
	@$(MAKE) --no-print-directory -C third_party/cplex