# Build the application
all:
    build
    test

# Build the application
build:
    echo "Building..."
    encore build

# Simulate an push on the CI
ci-push:
  act -s ENCORE_AUTH_KEY=$(echo $ENCORE_AUTH_KEY) -W .github/workflows/frontend.yaml push

# Run the application
run:
    encore run

dev:
    encore run --debug=break

# Test the application
test:
    gotestsum --raw-command -- encore test ./... -json

# Test the application
test-watch:
    gotestsum --watch --raw-command -- encore test ./... -json
