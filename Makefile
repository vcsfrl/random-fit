envfile := .env
include $(envfile)
export $(shell sed 's/=.*//' $(envfile))

# HELP
.PHONY: help

help: ## Usage: make <option>
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## APP Build.
	if [ ! -f .env ]; then cp -n .env.dist .env; echo "CONTAINER_EXEC_USER_ID=`id -u`" >> .env; echo "CONTAINER_USERNAME=${USER}" >> .env; fi
	docker compose build;

bash: ## APP Bash.
	docker compose run --remove-orphans random-fit_app bash

test: ## APP Test
	docker compose run --remove-orphans random-fit_app go test -race -cpu 24 -cover -coverprofile=data/test/coverage.out ./internal/...;

test-name: ##  Run test by name.
	docker compose run --remove-orphans random-fit_app go test -v -race -cpu 24 github.com/vcsfrl/random-fit/$(testPath) -run ^$(testName)$$;

test-debug:
	docker compose run  --remove-orphans --build --rm --service-ports random-fit_app /go/bin/dlv --listen=:$(RF_DEBUGGER_TEST_PORT) --headless=true --log=true --log-output=debugger,debuglineerr,gdbwire,lldbout,rpc --api-version=2 --accept-multiclient test  github.com/vcsfrl/random-fit/$(testPath) -- -test.run ^$(testName)$$;

lint: ## Run linter.
	docker run -t --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v2.0.2 golangci-lint run
