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
  act --var VITE_ENCORE_ENVIRONMENT=local --var ACTIONS_RUNTIME_TOKEN=foobar -s ENCORE_AUTH_KEY=$(echo $ENCORE_AUTH_KEY) push

# Simulate an push on the CI (Frontend)
ci-push-frontend:
  act --var VITE_ENCORE_ENVIRONMENT=local --var ACTIONS_RUNTIME_TOKEN=foobar -s ENCORE_AUTH_KEY=$(echo $ENCORE_AUTH_KEY) -W .github/workflows/frontend.yaml push

# Simulate an push on the CI (Backend)
ci-push-backend:
  act --var VITE_ENCORE_ENVIRONMENT=local --var ACTIONS_RUNTIME_TOKEN=foobar -s ENCORE_AUTH_KEY=$(echo $ENCORE_AUTH_KEY) -W .github/workflows/backend.yaml push

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
