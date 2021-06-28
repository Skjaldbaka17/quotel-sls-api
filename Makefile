.PHONY: build

build:
	sam build

run-local-with-env:
	sam local start-api --env-vars env.json
