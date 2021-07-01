.PHONY: build

build:
	sam build

build-run:
	sam build && make run-env

run-env:
	sam local start-api --env-vars env.json
	
check-swagger:
	which swagger || (go get -u github.com/go-swagger/go-swagger/cmd/swagger)
docs: check-swagger
	swagger generate spec --input ./swagger/tags.yaml -o ./swagger/swagger.yaml --scan-models -w ./functions
serve-docs:check-swagger
	swagger serve -F=swagger ./swagger/swagger.yaml
