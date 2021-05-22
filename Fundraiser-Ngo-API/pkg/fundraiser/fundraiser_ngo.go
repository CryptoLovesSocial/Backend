package fundraiser

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	ErrorFailedToUnmarshalRecord = "failed to unmarshal record"
	ErrorFailedToFetchRecord     = "failed to fetch record"
	ErrorInvalidUserData         = "invalid user data"
	ErrorCouldNotMarshalItem     = "could not marshal item"
	ErrorCouldNotDeleteItem      = "could not delete item"
	ErrorCouldNotDynamoPutItem   = "could not dynamo put item error"
	ErrorUserAlreadyExists       = "user.User already exists"
	ErrorUserDoesNotExists       = "user.User does not exist"
)

type FundraiserNgo struct {
	NgoId                  string `json:"pk"`
	FundraiserId           string `json:"sk"`
	FundraiserTitle        string `json:"fundraiserTitle"`
	FundraiserCause        string `json:"fundraiserCause"`
	FundraiserLocation     string `json:"fundraiserLocation"`
	FundraiserDescription  string `json:"fundraiserDescription"`
	FundraiserPhoto        string `json:"fundraiserPhoto"`
	FundraiserTargetAmount string `json:"fundraiserTargetAmount"`
}

func FetchFundraiserNgo(ngoId string, fundraiserId string, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*FundraiserNgo, error) {
	//Modifying the key for DynamoDB
	ngoId = "Ngo" + ngoId
	fundraiserId = "Fundraiser" + fundraiserId

	//Macking Call for DynamoDB
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String(ngoId),
			},
			"sk": {
				S: aws.String(fundraiserId),
			},
		},
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.GetItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)

	}

	//Sending the Get Request
	item := new(FundraiserNgo)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}
	return item, nil
}
func FetchFundraisersNgo(ngoId string, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*[]FundraiserNgo, error) {
	//Modifying the key for DynamoDB Storage
	ngoId = "Ngo" + ngoId

	//Macking Call for DynamoDB
	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {
				S: aws.String(ngoId),
			},
			":sk": {
				S: aws.String("Fundraiser"),
			},
		},
		KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :sk)"),
		TableName:              aws.String(tableName),
	}

	result, err := dynaClient.Query(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	var items *[]FundraiserNgo
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}
	return items, nil
}
func CreateFundraiserNgo(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (
	*FundraiserNgo,
	error,
) {
	//Checking if the correct request
	var u FundraiserNgo
	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, errors.New(ErrorInvalidUserData)
	}

	//Modifying the key for DynamoDB Storage
	u.NgoId = "Ngo" + u.NgoId
	u.FundraiserId = "Fundraiser" + u.FundraiserId

	//Marshaling the data
	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}
	//Puting it to DynamoDB
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}
	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}
	return &u, nil
}

func UpdateFundraiserNgo(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (
	*FundraiserNgo,
	error,
) {
	var u FundraiserNgo
	//Checking if the correct request
	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, errors.New(ErrorInvalidUserData)
	}

	// Check if Fundraiser exists
	currentFundraiser, _ := FetchFundraiserNgo(u.NgoId, u.FundraiserId, tableName, dynaClient)
	if currentFundraiser != nil && len(currentFundraiser.FundraiserId) == 0 {
		return nil, errors.New(ErrorUserDoesNotExists)
	}
	// Saveing it DynamoDB
	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}
	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}
	return &u, nil
}

func DeleteFundraiserNgo(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) error {
	//ngoId and fundraiserId from req
	ngoId := req.QueryStringParameters["ngoId"]
	fundraiserId := req.QueryStringParameters["fundraiserId"]

	//Modifying the key for DynamoDB
	ngoId = "Ngo" + ngoId
	fundraiserId = "Fundraiser" + fundraiserId

	//Deleting the Fundraiser
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String(ngoId),
			},
			"sk": {
				S: aws.String(fundraiserId),
			},
		},
		TableName: aws.String(tableName),
	}
	_, err := dynaClient.DeleteItem(input)
	if err != nil {
		return errors.New(ErrorCouldNotDeleteItem)
	}

	return nil
}
