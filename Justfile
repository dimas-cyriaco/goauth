# Start all
up:
  process-compose up -D

# Stop all
down:
  process-compose down

# Attach to process compose
attach:
  process-compose attach

# Build docker images
build:
    encore build docker goauth

# Tail backend test logs
test-logs-backend:
  process-compose process logs test-backend -f

# Run the application
run:
    encore run

# Start frontend dev server
dev-frontend:
  pnpm --prefix developer_area/frontend dev

# Start backend dev server with debugger
dev-backend:
    encore run --debug=break

# Start backend and frontend dev servers
dev-all:
  #!/usr/bin/env -S parallel --shebang --ungroup --jobs {{ num_cpus() }}
  just dev-backend
  just dev-frontend

# Test the application
test:
    gotestsum --raw-command -- encore test ./... -json

# Test the application in watch mode
test-watch:
    gotestsum --watch --raw-command -- encore test ./... -json

# Simulate a push on the CI
ci-push:
  act --var VITE_ENCORE_ENVIRONMENT=local --var ACTIONS_RUNTIME_TOKEN=foobar -s ENCORE_AUTH_KEY=$(echo $ENCORE_AUTH_KEY) push

# Simulate a push on the CI (Frontend)
ci-push-frontend:
  act --var VITE_ENCORE_ENVIRONMENT=local --var ACTIONS_RUNTIME_TOKEN=foobar -s ENCORE_AUTH_KEY=$(echo $ENCORE_AUTH_KEY) -W .github/workflows/frontend.yaml push

# Simulate a push on the CI (Backend)
ci-push-backend:
  act --var VITE_ENCORE_ENVIRONMENT=local --var ACTIONS_RUNTIME_TOKEN=foobar -s ENCORE_AUTH_KEY=$(echo $ENCORE_AUTH_KEY) -W .github/workflows/backend.yaml push
