

- [ ] Add get QOD for topics to /quotes/qod route
- [ ] Get qod/aod and histories by getting the newest inserted qod not by using current_date in postgres!
- [ ] Set all .ID that refer to AuthorIds to AuthorId (i.e. seen from API)
- [ ] Add 2-lambda functions and schedule them to run once daily inserting the QOD, QODICE, AOD and AODICES (and for QODs/QODICEs have it get random quote from topics for quality quotes)
- [ ] Put new DB on AWS
- [ ] Put new version of API on AWS
- [ ] Create new swagger based on newest version on AWS
- [ ] Put new API on RapidAPI
- [ ] Make README for future me better -- also readme for setup-quotel-db and old quotel-api (not serverless) and crawler
- [ ] Make new/Better README for RAPIDAPI
- [ ] optimize queries (longTime in tests)


- [ ] Explain how the searchin works in README (i.e. first plainto_ts then check fuzzy search if user had )
- [ ] Always order returned Quotes/Authors in /search by Popularity also (i.e. count desc)
- [ ] Testing by mocking GORM: https://betterprogramming.pub/how-to-unit-test-a-gorm-application-with-sqlmock-97ee73e36526 


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


