package email

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
)

const CharSet = "UTF-8"

type Send func(subject, msg string) error

func SendEmail(svc *ses.SES, sendFrom, emailTo string) Send {
	return func(subject, msg string) error {
		input := &ses.SendEmailInput{
			Destination: &ses.Destination{
				CcAddresses: []*string{},
				ToAddresses: []*string{
					aws.String(emailTo),
				},
			},
			Message: &ses.Message{
				Body: &ses.Body{
					Html: &ses.Content{
						Charset: aws.String(CharSet),
						Data:    aws.String(msg),
					},
					Text: &ses.Content{
						Charset: aws.String(CharSet),
						Data:    aws.String(msg),
					},
				},
				Subject: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(subject),
				},
			},
			Source: aws.String(sendFrom),
		}

		_, err := svc.SendEmail(input)

		return err
	}
}
