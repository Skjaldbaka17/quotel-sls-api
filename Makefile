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
	GO111MODULE=off swagger generate spec --input ./swagger/tags.yaml -o ./swagger/swagger.yaml --scan-models
serve-docs:check-swagger
	swagger serve -F=swagger ./swagger/swagger.yaml

upload-docs:
	make docs
	aws s3 cp ./swagger s3://www.api.quotel-rest.com --recursive --exclude "*" --include "swagger.*" --include "index.html" --acl "public-read"
	echo "\nDocs hosted HERE:\n\n\t http://www.api.quotel-rest.com.s3-website-eu-west-1.amazonaws.com\n"


test-authors:
#/authors tests
	cd ./functions/authors/get-aod-history && go test
	cd ./functions/authors/get-author-of-the-day && go test
	cd ./functions/authors/get-authors-by-ids && go test
	cd ./functions/authors/get-authors-list && go test
	cd ./functions/authors/get-random-author && go test

test-quotes:
#/quotes tests
	cd ./functions/quotes/get-qod-history && go test
	cd ./functions/quotes/get-quote-of-the-day && go test
	cd ./functions/quotes/get-quotes && go test
	cd ./functions/quotes/get-quotes-list && go test
	cd ./functions/quotes/get-random-quote && go test

test-search:
#/search tests
	cd ./functions/search/search-authors-by-string && go test
	cd ./functions/search/search-by-string && go test
	cd ./functions/search/search-quotes-by-string && go test

test-topics:
#/topics tests
	cd ./functions/topics/get-topic && go test
	cd ./functions/topics/get-topics && go test

test-meta:
	cd ./functions/meta/list-languages-supported && go test
	cd ./functions/meta/list-professions && go test
	cd ./functions/meta/list-nationalities && go test

test-daily:
#Test cron jobs
	cd ./functions/daily/set-of-the-day && go test

test-all:
	make test-authors
	make test-quotes
	make test-quotes
	make test-search
	make test-topics
	make test-daily
	make test-meta
	