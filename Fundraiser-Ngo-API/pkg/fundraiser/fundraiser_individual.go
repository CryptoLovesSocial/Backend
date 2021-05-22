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

type FundraiserIndividual struct {
	IndividualEmailId                string `json:"pk"`
	IndividualFundraiserId           string `json:"sk"`
	IndividualFirstname              string `json:"firstname"`
	IndividualLastname               string `json:"lastname"`
	IndividualPhoneNo                string `json:"phoneNo"`
	IndividualFundraiserTitle        string `json:"fundraiserTitle"`
	IndividualFundraiserCause        string `json:"fundraiserCause"`
	IndividualFundraiserLocation     string `json:"fundraiserLocation"`
	IndividualFundraiserDescription  string `json:"fundraiserDescription"`
	IndividualFundraiserPhoto        string `json:"fundraiserPhoto"`
	IndividualFundraiserTargetAmount string `json:"fundraiserTargetAmount"`
}

func FetchFundraiserIndividual(emailId string, fundraiserId string, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*FundraiserIndividual, error) {
	//Modifying the key for DynamoDB
	emailId = "Individual" + emailId
	fundraiserId = "Fundraiser" + fundraiserId

	//Macking Call for DynamoDB
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String(emailId),
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
	item := new(FundraiserIndividual)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}
	return item, nil
}
func FetchFundraisersIndividual(emailId string, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*[]FundraiserIndividual, error) {
	//Modifying the key for DynamoDB Storage
	emailId = "Individual" + emailId

	//Macking Call for DynamoDB
	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {
				S: aws.String(emailId),
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

	//Sending the Get Request
	var items *[]FundraiserIndividual
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}
	return items, nil
}
func CreateFundraiserIndividual(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (
	*FundraiserIndividual,
	error,
) {
	//Checking if the correct request
	var u FundraiserIndividual
	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, errors.New(ErrorInvalidUserData)
	}

	//Modifying the key for DynamoDB Storage
	u.IndividualEmailId = "Individual" + u.IndividualEmailId
	u.IndividualFundraiserId = "Fundraiser" + u.IndividualFundraiserId

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

func UpdateFundraiserIndividual(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (
	*FundraiserIndividual,
	error,
) {
	var u FundraiserIndividual
	//Checking if the correct request
	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, errors.New(ErrorInvalidUserData)
	}

	// Check if Fundraiser exists
	currentFundraiser, _ := FetchFundraiserIndividual(u.IndividualEmailId, u.IndividualFundraiserId, tableName, dynaClient)
	if currentFundraiser != nil && len(currentFundraiser.IndividualFundraiserId) == 0 {
		return nil, errors.New(ErrorUserDoesNotExists)
	}
	u.IndividualEmailId = "Individual" + u.IndividualEmailId
	u.IndividualFundraiserId = "Fundraiser" + u.IndividualFundraiserId

	// Saving it to DynamoDB
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

func DeleteFundraiserIndividual(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) error {
	//emailId and fundraiserId from req
	emailId := req.QueryStringParameters["emailId"]
	fundraiserId := req.QueryStringParameters["fundraiserId"]

	//Modifying the key for DynamoDB
	emailId = "Individual" + emailId
	fundraiserId = "Fundraiser" + fundraiserId

	//Deleting the Fundraiser
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String(emailId),
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
