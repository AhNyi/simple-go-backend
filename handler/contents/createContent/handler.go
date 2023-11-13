package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ahnyi/goformatvalidationerror"
	"github.com/ahnyi/simple-go-backend/common/util"
	"github.com/ahnyi/simple-go-backend/db"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"gopkg.in/go-playground/validator.v9"
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
	content := db.Content{}
	err := json.Unmarshal([]byte(request.Body), &content)
	if err != nil {
		return response(
			http.StatusBadRequest,
			errorResponseBody(err.Error()),
		), nil
	}

	validate := validator.New()
	err = validate.Struct(content)

	if errs, hasError := err.(validator.ValidationErrors); hasError {
		errFormat := goformatvalidationerror.New(errs)
		errContent, _ := json.Marshal(errFormat)

		return response(
			http.StatusBadRequest,
			string(errContent),
		), nil
	}

	now := util.NowFunc().Format(time.RFC3339)

	content.ID = uuid.New().String()
	content.CreatedAt = now
	content.UpdatedAt = now
	content.ItemType = db.ItemType.Content
	content.PreviewKey = uuid.New().String()

	if content.PublishedStartTime != "" {
		_publishStartTime, _ := time.Parse(time.RFC3339, content.PublishedStartTime)
		content.PublishedStartTime = _publishStartTime.Truncate(time.Minute).Format(time.RFC3339)
	}

	if content.PublishedEndTime != "" {
		_publishEndTime, _ := time.Parse(time.RFC3339, content.PublishedEndTime)
		content.PublishedEndTime = _publishEndTime.Truncate(time.Minute).Format(time.RFC3339)
	}

	err = DynamoDB.PutContentWithContext(ctx, content)
	if err != nil {
		return response(
			http.StatusInternalServerError,
			errorResponseBody(err.Error()),
		), nil
	}

	res, _ := json.Marshal(content)

	return response(http.StatusOK, string(res)), nil
}
