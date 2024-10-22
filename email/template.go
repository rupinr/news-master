package email

import (
	"bytes"
	"html/template"
)

type EmailData struct {
	ActivationLink string
}

func GenerateHTML(emailData EmailData) (string, error) {
	htmlTemplate := `
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Activate Your Subscription</title>
        <style>
            body {
                font-family: Arial, sans-serif;
                background-color: #f4f4f4;
                color: #333333;
                margin: 0;
                padding: 20px;
            }
            .container {
                max-width: 600px;
                margin: 0 auto;
                padding: 20px;
                background-color: #ffffff;
                box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
            }
            h1 {
                color: #0288D1;
            }
            a {
                color: #ffffff;
                background-color: #0288D1;
                padding: 10px 20px;
                text-decoration: none;
                border-radius: 5px;
            }
        </style>
    </head>
    <body>
        <div class="container">
            <h1>Welcome to QuickBrew!</h1>
            <p>We're thrilled to have you here. Click the link below to activate your subscription and start receiving personalized news updates directly to your inbox.</p>
            <p><a href="{{.ActivationLink}}" target="_blank">Activate Subscription</a></p>
            <p>Stay tuned for the latest news and updates, delivered fresh and hot, just like your favorite brew.</p>
            <p>Cheers,<br>The QuickBrew Team</p>
        </div>
    </body>
    </html>
    `

	t, err := template.New("email").Parse(htmlTemplate)
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
