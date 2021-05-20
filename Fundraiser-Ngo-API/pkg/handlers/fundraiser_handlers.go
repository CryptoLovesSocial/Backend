package handlers

import (
	"aws-lambda-api/pkg/fundraiser"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var ErrorMethodNotAllowed = "method Not allowed"

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

func GetFundraiser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (
	*events.APIGatewayProxyResponse,
	error,
) {
	ngoId := req.QueryStringParameters["ngoId"]
	fundraiserId := req.QueryStringParameters["fundraiserId"]
	result, err := fundraiser.FetchFundraiser(ngoId, fundraiserId, tableName, dynaClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
	}
	return apiResponse(http.StatusOK, result)
}
func GetFundraisers(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (
	*events.APIGatewayProxyResponse,
	error,
) {
	ngoId := req.QueryStringParameters["ngoId"]
	result, err := fundraiser.FetchFundraisers(ngoId, tableName, dynaClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
	}
	return apiResponse(http.StatusOK, result)
}

func CreateFundraiser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (
	*events.APIGatewayProxyResponse,
	error,
) {
	result, err := fundraiser.CreateFundraiser(req, tableName, dynaClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return apiResponse(http.StatusCreated, result)
}

func UpdateFundraiser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (
	*events.APIGatewayProxyResponse,
	error,
) {
	result, err := fundraiser.UpdateFundraiser(req, tableName, dynaClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return apiResponse(http.StatusOK, result)
}

func DeleteFundraiser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (
	*events.APIGatewayProxyResponse,
	error,
) {
	err := fundraiser.DeleteFundraiser(req, tableName, dynaClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return apiResponse(http.StatusOK, nil)
}

func UnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return apiResponse(http.StatusMethodNotAllowed, ErrorMethodNotAllowed)
}
