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

dev:
  encore run --debug=break

# Test the application
test:
    gotestsum --raw-command -- encore test ./... -json --debug

# Test the application
test-watch:
    gotestsum --watch --raw-command -- encore test ./... -json --debug
