package handler

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"

	"github.com/adzfaulkner/falconim/internal/email"
	"github.com/adzfaulkner/falconim/internal/recaptcha"
	"github.com/aws/aws-lambda-go/events"
)

type requestBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Subject  string `json:"subject"`
	Message  string `json:"message"`
	Response string `json:"response"`
}

type responseBody struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type logger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
}

func Handler(response CorsResponse, sendEmail email.Send, recaptchaChecker recaptcha.Checker, logger logger) func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		var reqBody requestBody
		err := json.Unmarshal([]byte(request.Body), &reqBody)

		if err != nil {
			logger.Error("Unable to unmarshal request body", zap.String("reqBody", request.Body), zap.Error(err))
			return *response(buildResponse(false, "Unable to process the request"), http.StatusBadRequest, nil), nil
		}

		if reqBody.Name == "" {
			return *response(buildResponse(false, "Name is required"), http.StatusBadRequest, nil), nil
		}

		if reqBody.Email == "" {
			return *response(buildResponse(false, "Email is required"), http.StatusBadRequest, nil), nil
		}

		if reqBody.Subject == "" {
			return *response(buildResponse(false, "Subject is required"), http.StatusBadRequest, nil), nil
		}

		if reqBody.Message == "" {
			return *response(buildResponse(false, "Message is required"), http.StatusBadRequest, nil), nil
		}

		if reqBody.Response == "" {
			return *response(buildResponse(false, "Response is required"), http.StatusBadRequest, nil), nil
		}

		reqIp := request.Headers["X-Forwarded-For"]

		res := recaptchaChecker(reqIp, reqBody.Response)

		if !res {
			logger.Error("Failed recaptcha check", zap.String("reqIp", reqIp), zap.String("response", reqBody.Response))
			return *response(buildResponse(false, "Unable to process the request"), http.StatusBadRequest, nil), nil
		}

		msg := fmt.Sprintf("<p>Name: %s</p><p>Email: %s</p><p>Message: %s</p>", reqBody.Name, reqBody.Email, reqBody.Message)

		err = sendEmail(reqBody.Subject, msg)

		if err != nil {
			logger.Error("Unable to send email", zap.Error(err))
			return *response(buildResponse(false, "Unable to process the request"), http.StatusBadRequest, nil), nil
		}

		return *response(buildResponse(true, "Email sent"), http.StatusCreated, nil), nil
	}
}

func buildResponse(success bool, msg string) []byte {
	rb := responseBody{
		Success: success,
		Message: msg,
	}

	res, _ := json.Marshal(&rb)

	return res
}
