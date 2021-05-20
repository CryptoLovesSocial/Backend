package ngo

import (
	"encoding/json"
	"errors"
	"log"

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

type Ngo struct {
	NgoId          string `json:"pk"`
	SK             string `json:"sk"`
	NgoName        string `json:"ngoName"`
	NgoAdress      string `json:"ngoAdress"`
	NgoCountry     string `json:"ngoCountry"`
	NgoDescription string `json:"ngoDescription"`
	NgoPhoto       string `json:"ngoPhoto"`
	NgoCategory    string `json:"ngoCategory"`
}

func FetchNgo(ngoId string, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*Ngo, error) {
	//Modifying the key for DynamoDB Storage
	ngoId = "Ngo" + ngoId

	//Macking Call for DynamoDB
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String(ngoId),
			},
			"sk": {
				S: aws.String("Details"),
			},
		},
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.GetItem(input)
	if err != nil {
		log.Println("err1")
		log.Println(err)
		return nil, errors.New(ErrorFailedToFetchRecord)

	}

	//Sending the Get Request
	item := new(Ngo)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		log.Println("err2")
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}
	return item, nil
}
func CreateNgo(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (
	*Ngo,
	error,
) {
	//Checking if the correct request
	var u Ngo
	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		println("err1")
		println(err)
		return nil, errors.New(ErrorInvalidUserData)
	}

	//Modifying the key for DynamoDB Storage
	u.NgoId = "Ngo" + u.NgoId
	u.SK = "Details"
	//Marshaling the data
	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		println("err2")
		println(err)
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}
	//Puting it to DynamoDB
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}
	_, err = dynaClient.PutItem(input)
	if err != nil {
		println("err3")
		println(err)
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}
	return &u, nil
}

func UpdateNgo(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (
	*Ngo,
	error,
) {
	var u Ngo
	//Checking if the correct request
	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, errors.New(ErrorInvalidUserData)
	}
	// Check if ngo exists
	currentNgo, _ := FetchNgo(u.NgoId, tableName, dynaClient)
	if currentNgo != nil && len(currentNgo.NgoId) == 0 {
		return nil, errors.New(ErrorUserDoesNotExists)
	}
	// Save ngo
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

func DeleteNgo(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) error {
	//ngoId from req
	ngoId := req.QueryStringParameters["ngoId"]

	//Modifying the key for DynamoDB Storage
	ngoId = "Ngo" + ngoId

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String(ngoId),
			},
			"sk": {
				S: aws.String("Details"),
			},
		},
		TableName: aws.String(tableName),
	}
	//Deleting the NGO
	_, err := dynaClient.DeleteItem(input)
	if err != nil {
		return errors.New(ErrorCouldNotDeleteItem)
	}

	return nil
}
