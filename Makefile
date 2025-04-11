.PHONY: build clean deps

# Install dependencies
deps:
	go mod tidy 

# Build the application
build:
	mkdir -p dist
	go build -o dist/dicedb-mcp

# Install the application
install:
	go install

# Clean build artifacts
clean:
	rm -rf dist/dicedb-mcp
