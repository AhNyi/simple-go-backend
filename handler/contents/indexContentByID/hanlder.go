package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ahnyi/simple-go-backend/db"
	"github.com/aws/aws-lambda-go/events"
)

var DynamoDB db.DB

func response(code int, body string) events.APIGatewayProxyResponse {
	defaultHeaders := map[string]string{"Content-Type": "application/json", "Access-Control-Allow-Origin": "*", "Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS", "Access-Control-Allow-Headers": "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization", "X-Content-Type-Options": "nosniff", "X-Frame-Options": "DENY"}
	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Body:       body,
		Headers:    defaultHeaders,
	}
}

func errorResponseBody(msg string) string {
	return fmt.Sprintf("{\"message\": \"%s\"}", msg)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	params := request.PathParameters
	if params["contractId"] == "" {
		return response(
			http.StatusBadRequest,
			errorResponseBody("ContractId is Required."),
		), nil
	}

	if params["contentId"] == "" {
		return response(
			http.StatusBadRequest,
			errorResponseBody("ContentId is Required."),
		), nil
	}

	content, err := DynamoDB.GetContentByIDWithContext(ctx, params["contentId"])
	if err != nil {
		return response(
			http.StatusBadRequest,
			errorResponseBody(err.Error()),
		), nil
	}

	res, _ := json.Marshal(content)
	return response(http.StatusOK, string(res)), nil
}
