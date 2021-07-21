# quotel-sls-api

This is the readme for the Quotel Serverless API hosted on AWS-lambda with the following endpoint 

https://rxvglshzhl.execute-api.eu-west-1.amazonaws.com/v1. 

Note you need an authorized API-key to access the API through that url. You can contact me at skjaldbaka17@gmail.com or https://www.linkedin.com/in/þórður-ágústsson/ to get an API key or you can use the rapidapi.com marketplace: https://rapidapi.com/skjaldbaka17/api/quotes-rest/ to get access to the API.

## About the API

The API is a connection to the Database setup by https://github.com/skjaldbaka17/setup-quotel-db which contains around 1.000.000 quotes, over 30.000 authors and over 100.000 quotes sorted into 133 topics. For documentation of this API see: http://www.api.quotel-rest.com.s3-website-eu-west-1.amazonaws.com .

The primary motivation for this project was to learn to setup and manage SaaS on AWS, using the serverless solutions provided by AWS like Lambda functions, API-Gateway and SAM. The project was a success.

## Requirements

* AWS CLI already configured with Administrator permission
* [Docker installed](https://www.docker.com/community-edition)
* [Golang](https://golang.org)
* SAM CLI - [Install the SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)

## Setup process

### Installing dependencies & building the target 

In this example we use the built-in `sam build` to automatically download all the dependencies and package our build target.   
Read more about [SAM Build here](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-cli-command-reference-sam-build.html) 

The `sam build` command is wrapped inside of the `Makefile`. To execute this simply run
 
```shell
make
```

## Local development

**Invoking function locally through local API Gateway**

Note you need to add a env.json file to root with "DATABASE_URL":"YOUR_DATABASE_URL" for the API to work, also remember that sam is running the API using Docker so if your DB is on your local machine you need to use as host:"docker.for.mac.localhost".

```bash
sam local start-api --env-vars env.json
```

You can also run

```bash
make build-run
```

which builds the project and then runs it locally using the env.json you created and put in root.

If one of the the previous commands ran successfully you should now be able to hit the following base local endpoint to invoke the lambda functions `http://localhost:3000/` (test with `http://localhost:3000/quotes/random` for getting a random quote)

### Testing

We use the `testing` package that comes built-in in Golang. Before running the tests you need to create a `.env` file in root with `DATABASE_URL=YOUR_DB_URL` then you can simply run the following commands to test the various functionalities of the api:

For /authors
```shell
make test-authors
```

For /quotes
```shell
make test-quotes
```
For /search
```shell
make test-search
```

For /topics
```shell
make test-topics
```

For /meta
```shell
make test-meta
```

For the daily cron jobs:
```shell
make test-daily
```

Or if you want to run all the above tests then run:

```shell
make test-all
```


### API Documentation

http://www.api.quotel-rest.com.s3-website-eu-west-1.amazonaws.com

For documenting the API we use Swagger (or OpenAPI) and document each endpoint inside the code with specific comments forexed with `swagger:route`. To compile these comments into a swagger.yaml file you simply run:

```shell
make docs
```

This command will first check if you have the goswagger bin compiled. If it is not installed on your machine the command should install it with the command
```shell
go get -u github.com/go-swagger/go-swagger/cmd/swagger
```

Then to upload the docs to the already made s3 bucket `s3://www.api.quotel-rest.com` (it compiles the comments first then uploads) you can simply run 

```shell
make upload-docs    
```

### The search

The search was implemented with phrases in mind. We wanted to try and make a fast full-text search for the quotes, which was a median success. I give the search about a 6.5/10. It does well on phrases but is extremely slow when searching for common single words (like 'love').

The general search works like this:

1. First we do a general phrase-search using `plainto_tsquery('english', 'love') as plainq` and then ordering the search based on `ts_rank(tsv, plainq)`, the ordering is a part of the reason why this takes so much time when the number of rows that the `( tsv @@ plainq )` returns is high. If this query returns some rows we return them to the user, otherwise we go to step 2.
2. Now the 1. search has 'failed' so we assume that the user has made a spelling error. FUZZY SEARCH!. We take their searchString and split it up into the distinct words (split based on space) and search for each word in our `unique_lexeme`, a materialized view containing each distinct stemmed word in our database (both from authors.name and quotes.quote), using a similarity search based on word trigrams. This will "fix" most common spelling error or atleast find the most similar word used in our database. For each word we find and fix we put them back into the sentence (searchString) and use this new (spelling fixed) string to search in the same way as in step 1. If this search returns some rows we return them to the user. If not we go to step 3.
3. Now most things have failed so we assume that the string the user sent us is a strange string (like trying to search for `Friedrich Nietzsche`, who knows if this is even a correct spelling of the man's name) and we assume the user is looking for some foreign author, therefore we take their string and do a similarity search against the `authors` table and if we get some matches we take the best match and find his/hers quotes and return those. If no match we return an empty array `[]`. 
4. Coming maybe later... Search for any quote/author containing at least one of the words in the searchString (not capable of implementing this now because the query takes tooooooo long a time on the small RDS machine I am using, maybe when I switch to using Aurora Postgres serverless service).

The search for quotes works like this:

1. Same as steps 1-2 in general-search but instead of running `ts_rank(tsv, plainq)` and `( tsv @@ plainq )` we run `ts_rank(quote_tsv, plainq)` and `( quote_tsv @@ plainq )` and also in step two we look only through distinc words of quotes, i.e. in materialized view `unique_lexeme_quotes`.

The search for authors works like this:

1. Same as steps 1-2 in general-search but instead of running `ts_rank(tsv, plainq)` and `( tsv @@ plainq )` we run `ts_rank(name_tsv, plainq)` and `( name_tsv @@ plainq )` and we run the queries against the `authors` table. And also in step two we look only through distinc words of authors names, i.e. in materialized view `unique_lexeme_authors`. If some rows are found we return them otherwise go to step 2.
2. Now we do a similarity search against the `authors.name`and return those rows who are most similar.


## Packaging and deployment

To deploy your application for the first time, run the following in your shell:

```bash
sam deploy --guided
```

The command will package and deploy your application to AWS, with a series of prompts:

You can find your API Gateway Endpoint URL in the output values displayed after deployment.


# Appendix

### Golang installation

Please ensure Go 1.x (where 'x' is the latest version) is installed as per the instructions on the official golang website: https://golang.org/doc/install

A quickstart way would be to use Homebrew, chocolatey or your linux package manager.

#### Homebrew (Mac)

Issue the following command from the terminal:

```shell
brew install golang
```

If it's already installed, run the following command to ensure it's the latest version:

```shell
brew update
brew upgrade golang
```

#### Chocolatey (Windows)

Issue the following command from the powershell:

```shell
choco install golang
```

If it's already installed, run the following command to ensure it's the latest version:

```shell
choco upgrade golang
```
