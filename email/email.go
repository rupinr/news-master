package email

import (
	"fmt"
	"log/slog"
	"news-master/app"
	"news-master/logger"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

var once sync.Once
var svc *ses.SES

func setupSession() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)
	if err != nil {
		slog.Error("Error in email")
	}

	svc = ses.New(sess)
}

const (
	CharSet = "UTF-8"
)

func sendSimulationEmail(recipient string, subject string, htmlBody string, textBody string) error {
	logger.Log.Debug(fmt.Sprintf("!!!Simulation!!!! Email sent to %v with text %v , html %v and subject %v", recipient, textBody, htmlBody, subject))
	return nil
}

func sendSesEmail(recipient string, subject string, htmlBody string, textBody string) error {

	once.Do(setupSession)
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(htmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(textBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(app.Config.EmailSender),
	}

	_, emailErr := svc.SendEmail(input)
	if emailErr == nil {
		logger.Log.Debug(fmt.Sprintf("Email sent to %v", recipient))
	}
	return emailErr

}

func SendEmail(recipient string, subject string, htmlBody string, textBody string) error {
	if app.Config.EmailSimulatorMode == string(OFF) {
		return sendSesEmail(recipient, subject, htmlBody, textBody)
	} else {
		return sendSimulationEmail(recipient, subject, htmlBody, textBody)
	}
}

type EmailMode string

const (
	ON  EmailMode = "on"
	OFF EmailMode = "off"
)
