package main

import (
	"github.com/adzfaulkner/falconim/internal/logger"
	"github.com/aws/aws-sdk-go/service/ses"
	"os"

	"github.com/adzfaulkner/falconim/cmd/internal/handler"
	"github.com/adzfaulkner/falconim/internal/email"
	"github.com/adzfaulkner/falconim/internal/recaptcha"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func main() {
	endpoint := os.Getenv("AWS_ENDPOINT")
	region := os.Getenv("AWS_REGION")
	corsAllowedOrigin := os.Getenv("CORS_ALLOWED_ORIGIN")

	awsCfg := aws.Config{Region: aws.String(region)}

	if endpoint != "" {
		awsCfg.Endpoint = aws.String(endpoint)
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            awsCfg,
		SharedConfigState: session.SharedConfigEnable,
	})

	if err != nil {
		panic(err.Error())
	}

	ssmsvc := ssm.New(sess, aws.NewConfig().WithRegion(region))
	sessvc := ses.New(sess, aws.NewConfig().WithRegion(region))

	emailFrom := getSmmParamVal(ssmsvc, "/FalconIM/EMAIL_FROM")
	emailTo := getSmmParamVal(ssmsvc, "/FalconIM/EMAIL_TO")
	recaptchaSecret := getSmmParamVal(ssmsvc, "/FalconIM/RECAPTCHA_SECRET")

	recaptchaCheck := recaptcha.Check(recaptchaSecret)
	sendEmail := email.SendEmail(sessvc, emailFrom, emailTo)

	lh, err := logger.NewHandler()

	if err != nil {
		panic(err.Error())
	}

	lambda.Start(handler.Handler(handler.GenerateResponse(corsAllowedOrigin), sendEmail, recaptchaCheck, lh))
}

func getSmmParamVal(ssmsvc *ssm.SSM, key string) string {
	param, err := ssmsvc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(true),
	})

	if err != nil {
		panic(err.Error())
	}

	return *param.Parameter.Value
}
