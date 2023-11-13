package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ahnyi/simple-go-backend/db"
	"github.com/aws/aws-lambda-go/events"
)

var DynamoDB db.DB

func response(code int, body string, headers ...map[string]string) events.APIGatewayProxyResponse {
	defaultHeaders := map[string]string{"Content-Type": "application/json", "Access-Control-Allow-Origin": "*", "Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS", "Access-Control-Allow-Headers": "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"}
	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Body:       body,
		Headers:    defaultHeaders,
	}
}

func errorResponseBody(msg string) string {
	return fmt.Sprintf("{\"message\":\"%s\"}", msg)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	return response(http.StatusOK, ""), nil
}
