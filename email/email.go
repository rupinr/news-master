package email

import (
	"fmt"
	"log/slog"
	"news-master/app"
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

func sendSimulationEmail(recipient string, subject string, htmlBody string, textBody string) {

	slog.Info(fmt.Sprintf("!!!Simulation!!!! Email sent to %v with text %v , html %v and subject %v", recipient, textBody, htmlBody, subject))
}

func sendSesEmail(recipient string, subject string, htmlBody string, textBody string) {

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

	_, err := svc.SendEmail(input)

	if err != nil {
		slog.Error(err.Error())
		return
	}
	slog.Debug("Email Sent to address: " + recipient)
}

func SendEmail(recipient string, subject string, htmlBody string, textBody string) {

	fmt.Println("env" + app.Config.EmailSimulatorMode)
	fmt.Println("struc" + string(OFF))

	if app.Config.EmailSimulatorMode == string(OFF) {
		sendSesEmail(recipient, subject, htmlBody, textBody)
	} else {
		sendSimulationEmail(recipient, subject, htmlBody, textBody)
	}
}

type EmailMode string

const (
	ON  EmailMode = "on"
	OFF EmailMode = "off"
)
