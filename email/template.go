package email

import (
	"bytes"
	"embed"
	"html/template"
)

type EmailData struct {
	ActivationLink string
}

//go:embed registration.html
var registrationTemplate embed.FS

func GenerateHTML(emailData EmailData) (string, error) {
	htmlTemplate, _ := registrationTemplate.ReadFile("registration.html")

	t, err := template.New("email").Parse(string(htmlTemplate))
	if err != nil {
		return "", err
	}

	var htmlBody bytes.Buffer
	if err = t.Execute(&htmlBody, emailData); err != nil {
		return "", err
	}

	return htmlBody.String(), nil
}

func GenerateText(emailData EmailData) (string, error) {
	textTemplate := `If the above link doesn't work, please copy and paste the following URL into your browser: {{.ActivationLink}}`

	t, err := template.New("email").Parse(textTemplate)
	if err != nil {
		return "", err
	}

	var textBody bytes.Buffer
	if err = t.Execute(&textBody, emailData); err != nil {
		return "", err
	}

	return textBody.String(), nil
}
