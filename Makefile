# Makefile for LogDrift

BINARY_NAME=logdrift
VERSION=0.1.0

.PHONY: all build test clean install run-test

all: build

build:
	@echo "\033[1;36mBuilding LogDrift v$(VERSION)...\033[0m"
	go build -o $(BINARY_NAME) *.go
	@echo "\033[1;32mBuild complete:\033[0m ./$(BINARY_NAME)"

test:
	@echo "Running tests..."
	go test -v ./...

run-test: build
	@echo "\033[1;36mTesting with sample logs...\033[0m"
	@chmod +x test_logs.sh
	./test_logs.sh | ./$(BINARY_NAME)

install: build
	@echo "\033[1;36mInstalling to /usr/local/bin...\033[0m"
	sudo cp $(BINARY_NAME) /usr/local/bin/
	@echo "\033[1;32mInstalled. Run: logdrift\033[0m"

clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	@echo "\033[1;32mClean complete\033[0m"

# Build for multiple platforms
build-all:
	@echo "\033[1;36mBuilding for multiple platforms...\033[0m"
	GOOS=darwin GOARCH=amd64 go build -o dist/$(BINARY_NAME)-darwin-amd64 *.go
	GOOS=darwin GOARCH=arm64 go build -o dist/$(BINARY_NAME)-darwin-arm64 *.go
	GOOS=linux GOARCH=amd64 go build -o dist/$(BINARY_NAME)-linux-amd64 *.go
	GOOS=linux GOARCH=arm64 go build -o dist/$(BINARY_NAME)-linux-arm64 *.go
	GOOS=windows GOARCH=amd64 go build -o dist/$(BINARY_NAME)-windows-amd64.exe *.go
	GOOS=windows GOARCH=arm64 go build -o dist/$(BINARY_NAME)-windows-arm64.exe *.go
	@echo "\033[1;32mMulti-platform build complete in ./dist/\033[0m"

# Quick development cycle
dev: build run-test

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Show help
help:
	@echo "LogDrift Build Commands:"
	@echo "  make build      - Build binary"
	@echo "  make run-test   - Build and test with sample logs"
	@echo "  make install    - Install to /usr/local/bin"
	@echo "  make clean      - Remove built binary"
	@echo "  make build-all  - Build for all platforms"
	@echo "  make dev        - Quick dev cycle (build + test)"
