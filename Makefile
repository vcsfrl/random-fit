envfile := .env
include $(envfile)
export $(shell sed 's/=.*//' $(envfile))

# Determine service name and run command based on APP_ENV
ifeq ($(APP_ENV),dev)
  SERVICE_NAME := random-fit_dev
  RUN_CMD := go run ./main.go
else
  SERVICE_NAME := random-fit
  RUN_CMD := ./bin/app
endif

COMPOSE_RUN := docker compose --profile $(APP_ENV) run --rm --remove-orphans

define require_dev
	@if [ "$(APP_ENV)" != "dev" ]; then echo "Error: '$(1)' requires APP_ENV=dev (current: $(APP_ENV))"; exit 1; fi
endef

# HELP
.PHONY: help

help: ## Usage: make <option>
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Build container image for the current APP_ENV (dev or prod).
	if [ ! -f .env ]; then cp -n .env.dist .env; echo "CONTAINER_EXEC_USER_ID=`id -u`" >> .env; echo "CONTAINER_USERNAME=${USER}" >> .env; fi
	docker compose --profile $(APP_ENV) build;

shell: ## Open a shell inside the container.
	docker compose --profile $(APP_ENV) down --remove-orphans
	docker compose --profile $(APP_ENV) up $(SERVICE_NAME) -d
	docker compose --profile $(APP_ENV) exec $(SERVICE_NAME) bash
	docker compose --profile $(APP_ENV) down --remove-orphans

generate: ## Generate code (dev only).
	$(call require_dev,generate)
	$(COMPOSE_RUN) $(SERVICE_NAME) go generate github.com/vcsfrl/random-fit
	$(COMPOSE_RUN) $(SERVICE_NAME) go generate github.com/vcsfrl/random-fit/cmd/translations

build: generate ## Build the binary (dev only).
	$(call require_dev,build)
	$(COMPOSE_RUN) $(SERVICE_NAME) go build -o ./bin/random-fit ./main.go

test: generate ## Run all tests (dev only).
	$(call require_dev,test)
	$(COMPOSE_RUN) $(SERVICE_NAME) go test -race -cpu 24 -cover -coverprofile=data/test/coverage.out ./...;

test-name: ## Run test by name (dev only).
	$(call require_dev,test-name)
	$(COMPOSE_RUN) $(SERVICE_NAME) go test -v -race -cpu 24 github.com/vcsfrl/random-fit/$(testPath) -run ^$(testName)$$;

test-debug: ## Debug a test (dev only).
	$(call require_dev,test-debug)
	docker compose --profile $(APP_ENV) run --remove-orphans --build --rm --service-ports $(SERVICE_NAME) /go/bin/dlv --listen=:$(RF_DEBUGGER_TEST_PORT) --headless=true --log=true --log-output=debugger,debuglineerr,gdbwire,lldbout,rpc --api-version=2 --accept-multiclient test  github.com/vcsfrl/random-fit/$(testPath) -- -test.run ^$(testName)$$;

lint: ## Run linter (dev only).
	$(call require_dev,lint)
	docker compose --profile $(APP_ENV) down --remove-orphans
	docker compose --profile $(APP_ENV) up $(SERVICE_NAME) -d
	docker compose --profile $(APP_ENV) exec $(SERVICE_NAME) /go/bin/golangci-lint run --timeout 5m
	docker compose --profile $(APP_ENV) down --remove-orphans

run: ## Run the app. Uses binary on prod, go run on dev. One-time container.
	$(COMPOSE_RUN) -i $(SERVICE_NAME) $(RUN_CMD)
