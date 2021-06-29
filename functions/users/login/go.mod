require (
	github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/structs v0.0.0-00010101000000-000000000000
	github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/utils v0.0.0-00010101000000-000000000000
	github.com/aws/aws-lambda-go v1.24.0
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e // indirect

)

replace github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/structs => ../../../local-dependencies/structs

replace github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/utils => ../../../local-dependencies/utils

module login

go 1.16
