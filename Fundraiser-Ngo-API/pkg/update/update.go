package update

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

type Update struct {
	FundraiserId      string `json:"pk"`
	UpdateId          string `json:"sk"`
	UpdateTitle       string `json:"updateTitle"`
	UpdateDescription string `json:"updateDescription"`
	UpdatePhoto       string `json:"updatePhoto"`
}

func FetchUpdate(fundraiserId string, updateId string, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*Update, error) {
	//Modifying the key for DynamoDB
	//FundraiserId as partition key and UpdateId as sort key
	fundraiserId = "Fundraiser" + fundraiserId
	updateId = "Update" + updateId

	//Macking Call for DynamoDB
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String(fundraiserId),
			},
			"sk": {
				S: aws.String(updateId),
			},
		},
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.GetItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)

	}

	//Sending the Get Request
	item := new(Update)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}
	return item, nil
}
func FetchUpdates(fundraiserId string, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*[]Update, error) {
	//Modifying the key for DynamoDB Storage
	fundraiserId = "Fundraiser" + fundraiserId

	//Macking Call for DynamoDB
	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {
				S: aws.String(fundraiserId),
			},
			":sk": {
				S: aws.String("Update"),
			},
		},
		KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :sk)"),
		TableName:              aws.String(tableName),
	}

	result, err := dynaClient.Query(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)

	}

	//Sending the Get Request
	var items *[]Update
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}
	return items, nil
}
func CreateUpdate(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (
	*Update,
	error,
) {
	//Checking if the correct request
	var u Update
	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, errors.New(ErrorInvalidUserData)
	}

	//Modifying the key for DynamoDB Storage
	u.FundraiserId= "Fundraiser" + u.FundraiserId
	u.UpdateId = "Update" + u.UpdateId

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

func UpdateUpdate(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (
	*Update,
	error,
) {
	var u Update
	//Checking if the correct request
	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, errors.New(ErrorInvalidUserData)
	}

	// Check if Update exists
	currentUpdate, _ := FetchUpdate(u.FundraiserId, u.UpdateId, tableName, dynaClient)
	if currentUpdate != nil && len(currentUpdate.UpdateId) == 0 {
		return nil, errors.New(ErrorUserDoesNotExists)
	}
	u.FundraiserId = "Fundraiser" + u.FundraiserId
	u.UpdateId = "Update" + u.UpdateId

	// Save Fundraiser
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

func DeleteUpdate(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) error {
	//fundraiserId and updateId from req
	fundraiserId := req.QueryStringParameters["fundraiserId"]
	updateId := req.QueryStringParameters["updateId"]

	//Modifying the key for DynamoDB
	fundraiserId = "Fundraiser" + fundraiserId
	updateId = "Update" + updateId

	//Deleting the Fundraiser
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String(fundraiserId),
			},
			"sk": {
				S: aws.String(updateId),
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
