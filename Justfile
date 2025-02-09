# Simple Justfile for a Go project

# Build the application
all:
    build
    test

# Build the application
build:
    echo "Building..."
    encore build

# Run the application
run:
    encore run

# Test the application
test:
    go run gotest.tools/gotestsum --raw-command -- encore test ./... -json

# Test the application
test-watch:
    go run gotest.tools/gotestsum --watch --raw-command -- encore test ./... -json
