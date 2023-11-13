package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-xray-sdk-go/xray"
)

var (
	Region                    = os.Getenv("REGION")
	TableName                 = os.Getenv("DB_NAME")
	DynamodbEndpoint          = os.Getenv("DYNAMODB_ENDPOINT")
	ContractID_ItemType_Index = "ContractID-ItemType-index"
	RowLockTimeout            = 180 * time.Second
	ItemType                  = TableItemType{
		Content: "Test#Content",
	}
)

type DB struct {
	Instance *dynamodb.DynamoDB
}

func New() DB {
	var sess *session.Session
	if val, ok := os.LookupEnv("ROWLOCK_TIMEOUT"); ok {
		if t, err := strconv.Atoi(val); err == nil {
			RowLockTimeout = time.Duration(t) * time.Second
		}
	}
	if DynamodbEndpoint != "" {
		sess = session.Must(session.NewSession(&aws.Config{
			Region:   aws.String(Region),
			Endpoint: aws.String(DynamodbEndpoint),
		}))
	} else {
		sess = session.Must(session.NewSession(&aws.Config{
			Region: aws.String(Region),
		}))
	}
	dynamo := dynamodb.New(sess)
	xray.AWS(dynamo.Client)

	return DB{Instance: dynamo}
}

//=============================================================================
// Common functions
//=============================================================================

func (d DB) Response(code int, body string, headers ...map[string]string) events.APIGatewayProxyResponse {
	defaultHeaders := map[string]string{"Content-Type": "application/json", "Access-Control-Allow-Origin": "*", "Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS", "Access-Control-Allow-Headers": "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization", "X-Content-Type-Options": "nosniff", "X-Frame-Options": "DENY"}
	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Body:       body,
		Headers:    defaultHeaders,
	}
}

func (d DB) ErrorResponseBody(msg string) string {
	return fmt.Sprintf("{\"message\":\"%s\"}", msg)
}

func (d DB) IndexByContractID(ctx context.Context, contractID string, itemType string) (*dynamodb.QueryOutput, error) {
	params := &dynamodb.QueryInput{
		TableName: aws.String(TableName),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":ContractID": {S: aws.String(contractID)},
			":ItemType":   {S: aws.String(itemType)},
		},
		KeyConditionExpression: aws.String("ContractID = :ContractID and ItemType = :ItemType"),
		IndexName:              aws.String(ContractID_ItemType_Index),
	}
	return d.Instance.QueryWithContext(ctx, params)
}

//=============================================================================
// Contents
//=============================================================================

func (d DB) PutContentWithContext(ctx context.Context, i interface{}) error {
	av, err := dynamodbattribute.MarshalMap(i)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(TableName),
	}

	_, err = d.Instance.PutItemWithContext(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

func (d DB) DeleteContentWithContext(ctx context.Context, contentId string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ID":       {S: aws.String(contentId)},
			"ItemType": {S: aws.String(ItemType.Content)},
		},
		TableName: aws.String(TableName),
	}

	_, err := d.Instance.DeleteItemWithContext(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

func (d DB) IndexContents(ctx context.Context, contractID string, itemType string) (contents []interface{}, err error) {
	res, err := d.IndexByContractID(ctx, contractID, itemType)
	if err != nil {
		return nil, err
	}
	if err = dynamodbattribute.UnmarshalListOfMaps(res.Items, &contents); err != nil {
		return nil, err
	}
	return contents, nil
}

func (d DB) GetContentByIDWithContext(ctx context.Context, contentId string) (content Content, err error) {
	params := &dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ID":       {S: aws.String(contentId)},
			"ItemType": {S: aws.String(ItemType.Content)},
		},
	}

	res, err := d.Instance.GetItemWithContext(ctx, params)
	if err != nil || len(res.Item) == 0 {
		return content, errors.New("GetContentById has an error or Item not found")
	}
	if err = dynamodbattribute.UnmarshalMap(res.Item, &content); err != nil {
		return content, err
	}

	return content, nil
}
