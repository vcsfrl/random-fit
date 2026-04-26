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

COMPOSE     := docker compose --profile $(APP_ENV)
COMPOSE_RUN := $(COMPOSE) run --rm --remove-orphans

define require_dev
	@if [ "$(APP_ENV)" != "dev" ]; then echo "Error: '$(1)' requires APP_ENV=dev (current: $(APP_ENV))"; exit 1; fi
endef

.DEFAULT_GOAL := help
.PHONY: build clean down generate help install lint run shell test test-debug test-name

help: ## Usage: make <option>
	@echo ""
	@echo "  APP_ENV = \033[33m$(APP_ENV)\033[0m  |  service = \033[33m$(SERVICE_NAME)\033[0m"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m  %-28s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""

install: ## Build container image for the current APP_ENV (dev or prod).
	@if [ ! -f .env ]; then cp -n .env.dist .env; echo "CONTAINER_EXEC_USER_ID=`id -u`" >> .env; echo "CONTAINER_USERNAME=${USER}" >> .env; fi
	$(COMPOSE) build

shell: ## Open a shell inside the container.
	$(COMPOSE) down --remove-orphans
	$(COMPOSE) up $(SERVICE_NAME) -d
	$(COMPOSE) exec $(SERVICE_NAME) bash
	$(COMPOSE) down --remove-orphans

ps: ## General: List running containers
	$(COMPOSE) ps;

generate: ## Generate code (dev only).
	$(call require_dev,generate)
	$(COMPOSE_RUN) $(SERVICE_NAME) sh -c 'go generate github.com/vcsfrl/random-fit && go generate github.com/vcsfrl/random-fit/cmd/translations'

build: generate ## Build the binary (dev only).
	$(call require_dev,build)
	$(COMPOSE_RUN) $(SERVICE_NAME) go build -o ./bin/random-fit ./main.go

test: generate ## Run all tests (dev only).
	$(call require_dev,test)
	$(COMPOSE_RUN) $(SERVICE_NAME) go test -race -cpu 24 -cover -coverprofile=data/test/coverage.out ./...

test-name: ## Run test by name (dev only).
	$(call require_dev,test-name)
	$(COMPOSE_RUN) $(SERVICE_NAME) go test -v -race -cpu 24 github.com/vcsfrl/random-fit/$(testPath) -run ^$(testName)$$

test-debug: ## Debug a test (dev only).
	$(call require_dev,test-debug)
	$(COMPOSE_RUN) --build --service-ports $(SERVICE_NAME) /go/bin/dlv --listen=:$(RF_DEBUGGER_TEST_PORT) --headless=true --log=true --log-output=debugger,debuglineerr,gdbwire,lldbout,rpc --api-version=2 --accept-multiclient test github.com/vcsfrl/random-fit/$(testPath) -- -test.run ^$(testName)$$

lint: ## Run linter (dev only).
	$(call require_dev,lint)
	docker run -t --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v2.114 golangci-lint run;

lint-fix: ## Dev: Run golangci-lint with auto-fix
	@$(call require-dev)
	docker run -t --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v2.11.4 golangci-lint run --fix;

run: ## Run the app. Uses binary on prod, go run on dev. One-time container.
	$(COMPOSE_RUN) --service-ports -i $(SERVICE_NAME) $(RUN_CMD) run

down: ## Stop and remove all containers for the current APP_ENV.
	$(COMPOSE) down --remove-orphans

clean: down ## Remove build artifacts and stop containers.
	@rm -rf ./bin/*
	@echo "Build artifacts cleaned."
