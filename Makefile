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

shell: ## APP Shell.
	docker compose down --remove-orphans
	docker compose up random-fit_dev -d
	docker compose exec random-fit_dev bash
	docker compose down --remove-orphans


generate: ## APP Generate code.
	docker compose run --remove-orphans random-fit_dev go generate github.com/vcsfrl/random-fit
	docker compose run --remove-orphans random-fit_dev go generate github.com/vcsfrl/random-fit/cmd/translations

build: generate## APP Build.
	docker compose run --remove-orphans random-fit_dev go build -o ./bin/random-fit ./main.go

test: generate ## APP Test
	docker compose run --remove-orphans random-fit_dev go test -race -cpu 24 -cover -coverprofile=data/test/coverage.out ./...;

test-name: ##  Run test by name.
	docker compose run --remove-orphans random-fit_dev go test -v -race -cpu 24 github.com/vcsfrl/random-fit/$(testPath) -run ^$(testName)$$;

test-debug:
	docker compose run  --remove-orphans --build --rm --service-ports random-fit_dev /go/bin/dlv --listen=:$(RF_DEBUGGER_TEST_PORT) --headless=true --log=true --log-output=debugger,debuglineerr,gdbwire,lldbout,rpc --api-version=2 --accept-multiclient test  github.com/vcsfrl/random-fit/$(testPath) -- -test.run ^$(testName)$$;

lint: ## Run linter.
	#docker run -t --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v2.0.2 golangci-lint run
	docker compose down --remove-orphans
	docker compose up random-fit_dev -d
	docker compose exec random-fit_dev /go/bin/golangci-lint run --timeout 5m
	docker compose down --remove-orphans

build-docker-image:
	docker build --build-arg username=rf --build-arg exec_user_id=1000  -t vcsfrl/random-fit:v1.0.0 --target prod .
	#docker run --rm -it --entrypoint bash vcsfrl/random-fit:v1.0.0
	#docker tag <image hash> <cluster ip>:32000/vcsfrl/random-fit:v1.0.0
	#docker push <cluster ip>:32000/vcsfrl/random-fit:v1.0.0

run: ## Run the app.
	docker compose run --remove-orphans random-fit
