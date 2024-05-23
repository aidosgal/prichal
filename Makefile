# Go parameters
GOCMD=CONFIG_PATH="./config/local.yaml" go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=app
BINARY_UNIX=$(BINARY_NAME)_unix

# Default target executed when you run `make`
all: run 

# Run the application
run:
	$(GOCMD) run cmd/prichal/main.go

# Build the application
build: 
	$(GOBUILD) -o $(BINARY_NAME) -v

# Install dependencies
deps:
	$(GOGET) -u ./...

# Test the application
test: 
	$(GOTEST) -v ./...

# Clean build artifacts
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Cross-compile for Linux
build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

# Include a help target to describe make targets
help:
	@echo "Makefile commands:"
	@echo "  all          - Runs tests and builds the application"
	@echo "  run          - Runs the application"
	@echo "  build        - Builds the application"
	@echo "  deps         - Installs dependencies"
	@echo "  test         - Runs tests"
	@echo "  clean        - Cleans build artifacts"
	@echo "  build-linux  - Cross-compiles the application for Linux"

