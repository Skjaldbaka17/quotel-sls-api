module github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/utils

go 1.16

replace github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/structs => ../structs

require (
	github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/structs v0.0.0-00010101000000-000000000000
	github.com/aws/aws-lambda-go v1.24.0
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.11
)
