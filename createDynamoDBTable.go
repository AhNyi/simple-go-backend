package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func createDynamoDBTable() error {
	// Create a new AWS session
	sess := session.Must(session.NewSession(&aws.Config{}))

	// Create a DynamoDB client using the session
	client := dynamodb.New(sess)

	// Define the parameters for createing a DynamoDB table
	input := &dynamodb.CreateTableInput{
		TableName: aws.String("simple-go-test"),
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       aws.String("HASH"), // Partition key
			},
			{
				AttributeName: aws.String("ItemType"),
				KeyType:       aws.String("RANGE"), // Sort key
			},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: aws.String("S"), // Assuming ID is a string
			},
			{
				AttributeName: aws.String("ItemType"),
				AttributeType: aws.String("S"), // Assuming ItemType is a string
			},
			{
				AttributeName: aws.String("ContractID"),
				AttributeType: aws.String("S"), // Assuming ContractID is a string
			},
		},
		BillingMode: aws.String("PROVISIONED"),
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("Contract-ItemType-index"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("ContractID"),
						KeyType:       aws.String("HASH"), // Index partition key
					},
					{
						AttributeName: aws.String("ItemType"),
						KeyType:       aws.String("RANGE"), // Index sort key
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(5),
					WriteCapacityUnits: aws.Int64(5),
				},
			},
		},
	}

	// Create the DynamoDB table
	_, err := client.CreateTable(input)
	if err != nil {
		return err
	}

	fmt.Println("DynamoDB table created successfully!")
	return nil
}

func main() {
	if err := createDynamoDBTable(); err != nil {
		fmt.Println("Error: ", err)
	}
}
