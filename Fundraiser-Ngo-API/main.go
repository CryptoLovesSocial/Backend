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
	//Handling request of NGO's
	//PartitionKey = constant string of DetailsNGO
	//SortKey = NgoId
	case "GET" + "|" + "getNgo":
		return handlers.GetNgo(req, tableName, dynaClient)
	case "POST" + "|" + "createNgo":
		return handlers.CreateNgo(req, tableName, dynaClient)
	case "PUT" + "|" + "updateNgo":
		return handlers.UpdateNgo(req, tableName, dynaClient)
	case "DELETE" + "|" + "deleteNgo":
		return handlers.DeleteNgo(req, tableName, dynaClient)
	case "GET" + "|" + "getNgos":
		return handlers.GetNgos(req, tableName, dynaClient)

	//Handling request of Fundraiser -> NGO(s)
	//PartitionKey = NgoId
	//SortKey = FundraiserId
	case "GET" + "|" + "getFundraiserNgo":
		return handlers.GetFundraiserNgo(req, tableName, dynaClient)
	case "POST" + "|" + "createFundraiserNgo":
		return handlers.CreateFundraiserNgo(req, tableName, dynaClient)
	case "PUT" + "|" + "updateFundraiserNgo":
		return handlers.UpdateFundraiserNgo(req, tableName, dynaClient)
	case "DELETE" + "|" + "deleteFundraiserNgo":
		return handlers.DeleteFundraiserNgo(req, tableName, dynaClient)
	case "GET" + "|" + "getFundraisersNgo":
		return handlers.GetFundraisersNgo(req, tableName, dynaClient)

	//Handling request of Fundraiser -> Individual(s)
	//PartitionKey = IndividualEmailId
	//SortKey = FundraiserId
	case "GET" + "|" + "getFundraiserIndividual":
		return handlers.GetFundraiserIndividual(req, tableName, dynaClient)
	case "POST" + "|" + "createFundraiserIndividual":
		return handlers.CreateFundraiserIndividual(req, tableName, dynaClient)
	case "PUT" + "|" + "updateFundraiserIndividual":
		return handlers.UpdateFundraiserIndividual(req, tableName, dynaClient)
	case "DELETE" + "|" + "deleteFundraiserIndividual":
		return handlers.DeleteFundraiserIndividual(req, tableName, dynaClient)
	case "GET" + "|" + "getFundraisersIndividual":
		return handlers.GetFundraisersIndividual(req, tableName, dynaClient)

	//Handling request of Update -> Fundraiser(s)
	//PartitionKey = FundraiserId
	//SortKey = UpdateId
	case "GET" + "|" + "getUpdate":
		return handlers.GetUpdate(req, tableName, dynaClient)
	case "POST" + "|" + "createUpdate":
		return handlers.CreateUpdate(req, tableName, dynaClient)
	case "PUT" + "|" + "updateUpdate":
		return handlers.UpdateUpdate(req, tableName, dynaClient)
	case "DELETE" + "|" + "deleteUpdate":
		return handlers.DeleteUpdate(req, tableName, dynaClient)
	case "GET" + "|" + "getUpdates":
		return handlers.GetUpdates(req, tableName, dynaClient)
	default:
		return handlers.UnhandledMethod()
	}
}
