package main

import (
	"github.com/ahnyi/simple-go-backend/db"
	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	DynamoDB = db.New()
}

func main() {
	lambda.Start(handler)
}