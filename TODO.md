
- [ ] Put new API on RapidAPI
- [ ] Make new/Better README for RAPIDAPI
- [ ] Use RapidAPI for WhoTheFuckSaidThat.com

- [ ] Make README for future me better -- also readme for setup-quotel-db and old quotel-api (not serverless) and crawler
- [ ] Explain how the searchin works in README (i.e. first plainto_ts then check fuzzy search if user had )

- [ ] Setup AURORA POSTGRES SERVERLESS with the new Data (pg_dump? https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/Aurora.Migrate.html -> https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/AuroraPostgreSQL.Migrating.html#AuroraPostgreSQL.Migrating.RDSPostgreSQL.Import.Console)






- [x] set utils/structs as packages in their own module called 'handlers' or something like that
- [x] Optimize random for both English and Icelandic and BOTH!!!
- [x] Developer portal: https://docs.amazonaws.cn/en_us/apigateway/latest/developerguide/apigateway-developer-portal.html !!!
- [x] https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-developer-portal.html#apigateway-developer-portal-create !!!!!
- [x] https://docs.amazonaws.cn/en_us/apigateway/latest/developerguide/rest-api-distribute.html !!!
- [x] https://serverlessrepo.aws.amazon.com/applications/arn:aws:serverlessrepo:us-east-1:563878140293:applications~api-gateway-dev-portal !!!! 
- [x] Authorizer instead of yourself?
- [x] Let ApiGateway / aws take care of Api_key authorization etc (https://aws.amazon.com/blogs/compute/generate-your-own-api-gateway-developer-portal/ --- usage plans)
- [x] AWS marketplace for payement of api etc?! (rapidapi.com)
- [x] Setup A testing structure and test all functions already created + TDD after that
- [x] Add swagger and docs route
- [x] Adapt to new data + add tests for all
- [x] Add route in /meta for getting all distinct Professions and distinct Nationalities
- [x] Make all tests use 'GetRequest'
- [x] Testing by mocking GORM: https://betterprogramming.pub/how-to-unit-test-a-gorm-application-with-sqlmock-97ee73e36526  #I want to test the REAL db, full text search test etc
- [x] Add get QOD for topics to /quotes/qod route
- [x] Get qod/aod and histories by getting the newest inserted qod not by using current_date in postgres! (not exactly, am not using this method as fail safe just in case)
- [x] Set all .ID that refer to AuthorIds to AuthorId (i.e. seen from API)
- [x] Put new DB on AWS (Image/copy)
- [x] Put new version of API on AWS
- [x] Add 2-lambda functions and schedule them to run once daily inserting the QOD, QODICE, AOD and AODICES (and for QODs/QODICEs have it get random quote from topics for quality quotes)
- [x] Create new swagger based on newest version on AWS (Multiple Examples?)
- [x] Make examples for all Features RapidAPI (Maybe in SWAGGER?) -- Did it in postman
- [x] optimize queries (longTime in tests)
- [x] Change 'deathDate' and 'birthDate' in api to 'died' and 'born'?
- [x] When: "request body is not structured correctly. Please refer to the /docs page for information on how to structure the request body" then set instead of '/docs' the link
to the S3 hosted website with the docs (even better let the link be specific to the endpoint that cause it for example http://www.api.quotel-rest.com.s3-website-eu-west-1.amazonaws.com/#operation/GetQuotes)


//Dumping local database to a file called pgexpdump.sql
pg_dump -U <USER_NAME> dbname=<DATABASE_NAME> -f <FILE_NAME_WHERE_TO_STORE_DUMP>.sql

//restoring the data to the db
psql \
	-f <FILE_NAME_WHERE_DUMP_IS_STORED>.sql \
	--host=<HOST_NAME> \
   --port=<PORT> \
   --username=<USER_NAME> \
   --password  \
   --dbname=<DB_NAME>
