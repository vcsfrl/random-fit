
# HELP
.PHONY: help

help: ## Usage: make <option>
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## APP Build.
	if [ ! -f .env ]; then cp -n .env.dist .env; echo "CONTAINER_EXEC_USER_ID=`id -u`" >> .env; echo "CONTAINER_USERNAME=${USER}" >> .env; fi
	docker compose build;

shell: ## APP Bash.
	docker compose run --remove-orphans random-fit_app bash

test: ## APP Test
	go test -v -race -cpu 24 -cover -coverprofile=data/test/coverage.out ./...;

test-name: ##  Run test by name.
	go test -v -race -cpu 24 github.com/vcsfrl/random-fit/$(testPath) -run ^$(testName)$$;

lint: ## Run linter.
	docker run -t --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v2.0.2 golangci-lint run
