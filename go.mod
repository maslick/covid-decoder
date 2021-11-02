module github.com/maslick/covid-decoder

go 1.16

require (
	github.com/aws/aws-lambda-go v1.27.0
	github.com/awslabs/aws-lambda-go-api-proxy v0.11.0
	github.com/stapelberg/coronaqr v0.0.0-20210811192332-c5f1f57c0e99
)

replace github.com/stapelberg/coronaqr => github.com/denysvitali/coronaqr v0.0.0-20211027233650-9c753cb0f96f
