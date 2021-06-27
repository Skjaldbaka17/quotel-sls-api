require (
	github.com/Skjaldbaka17/quotel-sls-api/structs v0.0.0-00010101000000-000000000000
	github.com/aws/aws-lambda-go v1.23.0
)

module get-quotes

go 1.16

replace github.com/Skjaldbaka17/quotel-sls-api/structs => ../../structs
