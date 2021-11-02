//go:build lambda
// +build lambda

package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/maslick/covid-decoder/src"
)

type LambdaController struct {
	src.RestController
}

func (ctrl *LambdaController) LambdaHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	router := ctrl.Init()
	return httpadapter.New(router).ProxyWithContext(ctx, req)
}

func main() {
	ctrl := LambdaController{src.RestController{Service: &src.Service{}}}
	lambda.Start(ctrl.LambdaHandler)
}
