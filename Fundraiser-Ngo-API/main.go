package main

import (
	"aws-lambda-api/pkg/handlers"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	dynaClient dynamodbiface.DynamoDBAPI
)

func main() {
	region := os.Getenv("AWS_REGION")
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return
	}
	dynaClient = dynamodb.New(awsSession)
	lambda.Start(handler)
}

const tableName = "NGOdetails"

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod + "|" + req.PathParameters["method"] {
	case "GET" + "|" + "getfundraiser":
		return handlers.GetFundraiser(req, tableName, dynaClient)
	case "POST" + "|" + "createfundraiser":
		return handlers.CreateFundraiser(req, tableName, dynaClient)
	case "PUT" + "|" + "updatefundraiser":
		return handlers.UpdateFundraiser(req, tableName, dynaClient)
	case "DELETE" + "|" + "deletefundraiser":
		return handlers.DeleteFundraiser(req, tableName, dynaClient)
	case "GET" + "|" + "getfundraisers":
		return handlers.GetFundraisers(req, tableName, dynaClient)

	case "GET" + "|" + "getngo":
		return handlers.GetNgo(req, tableName, dynaClient)
	case "POST" + "|" + "createngo":
		return handlers.CreateNgo(req, tableName, dynaClient)
	case "PUT" + "|" + "updatengo":
		return handlers.UpdateNgo(req, tableName, dynaClient)
	case "DELETE" + "|" + "deletengo":
		return handlers.DeleteNgo(req, tableName, dynaClient)
	default:
		return handlers.UnhandledMethod()
	}
}
