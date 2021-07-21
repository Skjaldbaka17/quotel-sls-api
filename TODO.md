# TODOs
- [ ] Put new API on RapidAPI
- [ ] Make new/Better README for RAPIDAPI
- [ ] Use RapidAPI for WhoTheFuckSaidThat.com

- [ ] Setup AURORA POSTGRES SERVERLESS with the new Data (pg_dump? https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/Aurora.Migrate.html -> https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/AuroraPostgreSQL.Migrating.html#AuroraPostgreSQL.Migrating.RDSPostgreSQL.Import.Console)


 ---------------------------- WHOTHEFUCKSAIDTHAT.COM ---------------------------- 
- [ ] Google Analytics for the site (set it up on google)
- [x] https certificate



 ---------------------------- CUSTOM FRONT END FOR API USERS ---------------------------- 

- [ ] Front end for API (Create it in its own repo!) (Find template?)
- [ ] LandingPage (with minor info i.e. used by www.whothefucksaidthat.com + some quotes + tiers/pricing info)
- [ ] SignUp / Login Using aws Cognito
- [ ] Move Users Backend to Front End Repo?
- [ ] HomePage for users (History of requests + Tier + upgrade / downgrade tier)
- [ ] Pay with Crypto



 ---------------------------- FURTHER STUFF ---------------------------- 

- [ ] Draw up DB-Graph (i.e. how tables are connected to view etc)

- [ ] Insert Quote for created author or for a 'real' author (private and public)
- [ ] update inserted quote (priv and pub)
- [ ] Create new Author (private and public)
- [ ] Update created author (priv and pub)
- [ ] Create new Topic (private and public)
- [ ] update created topic (priv and pub)
- [ ] Sort return list alphabetically Icelandic support

- [ ] Look into payment for some privileges




 ---------------------------- DONE SERVERLESS ---------------------------- 

- [x] New crawler for new quotes / authors
- [x] is random truly random (i.e. does the "random" funcitonality truly return randomly or is it biased towards quotes in the "front" of the DB (i.e. in the front where postgres stores them)) -- now using `tablesample system(0.1)`if whole table otherwise using `order by random()`
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
- [x] Explain how the searchin works in README (i.e. first plainto_ts then check fuzzy search if user had )
- [x] Make README for future me better -- also readme for setup-quotel-db and old quotel-api (not serverless) and crawler




 ---------------------------- DONE DEDICATED SERVER ---------------------------- 


- [x] Make Authors Search more efficient (create a similarity-based index ?)
- [x] Change to Use Gorm to the fullest, oooooooorrr just change returned json to : {"name":"authorName", "id":authorId, "hasIcelandicQuotes":true/false, "nrOfEnglishQuoes":int, "nrOfIcelandicQuotes":int, "quotes":[{"quote": "theQuote", "id":quoteId, "isIcelandic": true/false}]}
- [x] SEARCHSPEED VERY SLOW ON SERVER (10-15sec!) => indexes + move to DynamoDB/NOSql!
- [x] QOD/AOD history called -> if QOD/AOD not set for some day in the history then create that row for that date! -> Made a cron job in AWS lambda (setOfTheDay see https://github.com/skjaldbaka17/quotel-sls-api)
- [x] WebsiteLogo
- [x] automate setting up the EC2 and fetching the code and running server -- basically did that with aws lambda (see https://github.com/skjaldbaka17/quotel-sls-api )
- [x] Separate WhoTheFuckSaidThat.com from the API (i.e. have as its own APP that queries the API! + its own Repo)
- [x] Find cheapest option to run the API and the API/USER front end and use that!
- [x] Test/Use AWS lambda (Cheaper?)
- [x] Test/Use Amazon's API Gateway + Lambda (Cheaper?) //https://www.quora.com/What-is-the-best-and-cheapest-way-of-hosting-REST-API-as-a-startup-I-am-using-AWS-EC2-but-I-am-not-sure-whether-that-is-the-best-option-or-not-for-the-startup-who-has-limited-budget
- [x] Copy The PostgresDB (as it is just after setup) into an S3 bucket for safekeeping -- at least made pg_dump/pg_restore
- [x] Frontend Look fixes according to Roberto
- [x] More info about author (Wikipedia link + birth-death i.e. for example 1901-2000)
- [x] Buy and setup Domain name
- [x] Make Random query faster
- [x] Setup frontend on / route for displaying a random quote.
- [x] Setup AWS (EC2) server
- [x] Setup RDS-Postgres db on AWS and setup the quotes
- [x] Add .env variables (i.e. for names of tables etc.)
- [x] Coordinate naming convention (apiKey vs api_key vs apikey etc)
- [x] only return keys, in the response-json, that are relevant to the request
- [x] CleanUp DB after tests
- [x] add api key to swagger
- [x] Add authentication for access to the api + Creating apiKeys + Documenting usage + admin access vs normal access
- [x] Add password protection / protected routes capability (at least for SetQuoteOfTheyDay route )
- [x] Save History of requests
- [x] Add Users (GOD vs ...)
- [x] Error handling
- [x] Add error response to Swagger
- [x] Go over Swagger + Clean it up and make pretty
- [x] Clean up Documentation look (Swagger)
- [x] Review /topics for Swagger 
- [x] Review /search for Swagger 
- [x] Review /quotes for Swagger 
- [x] Review /meta for Swagger 
- [x] Review /authors for Swagger 
- [x] Clean up get/set QOD/AOD
- [x] Pagination Everywhere where needed
- [x] Clean up routes files
- [x] add search for topics and searching in topics (maybe just have a single search endpoint with parameters? i.e. want to search for authors / quotes / language inside a specific topic?)
- [x] getQuotes route (combine with getQuotesById and add to it to get quotes from a specific author + add pagination)
- [x] Add Counting each time a quote is accessed / sent from Api (also for topics) - i.e. stats
- [x] Add tests for GetAOD and AODHistory and SetAOD
- [x] Add tests for GetQOD and QODHistory and SetQOD
- [x] Add get and set Author of The Day (plus points for available to set authors for multiple days in one request + plus points for AOD history)
- [x] Add get and set Quote of The Day (plus points for available to set quotes for multiple days in one request + plus points for QOD history)
- [x] quoteoftheday << qod
- [x] Get list of authors (with parameters for pagination and alphabet and languages)
- [x] Get random Author
- [x] Validate RequestBody + Tests
- [x] Add get random Quote
- [x] Use TopicsView instead of searchview (and change name to something more general)
- [x] Add Icelandic / English Support
- [x] Add categories
- [x] English and Icelandic Authors with same name have same author id
- [x] Add Search-"scroll", User is searching and is scrolling through her search and wants next batch of results matching her search i.e. PAGINATION
- [x] setup testing (unit)
- [x] Implement GetQuotesById (multiple quotes route)
- [x] Clean-up test files (Move some lines into their own functions etc.)
- [x] Setup Swagger for api-docs 
      * https://github.com/go-swagger/go-swagger
      * https://github.com/nicholasjackson/building-microservices-youtube/blob/episode_7/product-api/main.go
      * https://www.youtube.com/watch?v=07XhTqE-j8k&t=374s
      * https://github.com/nicholasjackson/building-microservices-youtube/blob/episode_7/product-api/handlers/docs.go


Author api: https://quotes.rest




### Helpful resources for full-text search in postgres

* https://www.opsdash.com/blog/postgres-full-text-search-golang.html 
* https://medium.com/@bencagri/implementing-multi-table-full-text-search-with-gorm-632518257d15
* https://www.freecodecamp.org/news/fuzzy-string-matching-with-postgresql/
* https://www.compose.com/articles/mastering-postgresql-tools-full-text-search-and-phrase-search/ 



```bash
//Dumping local database to a file called pgexpdump.sql
pg_dump -U <USER_NAME> dbname=<DATABASE_NAME> -f <FILE_NAME_WHERE_TO_STORE_DUMP>.sql
```

```bash
//restoring the data to the db
psql \
	-f <FILE_NAME_WHERE_DUMP_IS_STORED>.sql \
	--host=<HOST_NAME> \
   --port=<PORT> \
   --username=<USER_NAME> \
   --password  \
   --dbname=<DB_NAME>
```
