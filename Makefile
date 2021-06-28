.PHONY: build

build:
	sam build

run-env:
	sam local start-api --env-vars env.json
