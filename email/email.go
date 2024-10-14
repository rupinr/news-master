package email

import "fmt"

func SendEmail(to string, body string, subject string) {

	fmt.Printf("Sending email to %v with token %v", to, body)
}
