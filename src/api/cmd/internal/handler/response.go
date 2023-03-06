package handler

import (
	"github.com/aws/aws-lambda-go/events"
)

type ResponseBody []byte
type ResponseStatus int
type ResponseHeaders map[string]string

type CorsResponse func(body ResponseBody, status ResponseStatus, headers ResponseHeaders) *events.APIGatewayProxyResponse

func GenerateResponse(allowOrigin string) CorsResponse {
	return func(body ResponseBody, status ResponseStatus, headers ResponseHeaders) *events.APIGatewayProxyResponse {
		if headers == nil {
			headers = map[string]string{}
		}

		headers["Content-Type"] = "application/json"
		headers["Access-Control-Allow-Origin"] = allowOrigin

		return &events.APIGatewayProxyResponse{
			Body:       string(body),
			StatusCode: int(status),
			Headers:    headers,
		}
	}
}
