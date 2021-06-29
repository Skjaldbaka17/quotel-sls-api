.PHONY: build

build:
	sam build

build-run:
	sam build && make run-env

run-env:
	sam local start-api --env-vars env.json
