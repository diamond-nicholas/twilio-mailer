package controllers

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendEmail(toEmail, toName, subject, templatePath string, data interface{}) error {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		return err
	}

	from := mail.NewEmail("Your App Name", os.Getenv("SENDGRID_FROM"))
	to := mail.NewEmail(toName, toEmail)
	plainTextContent := "This is a plain text version of the email."
	htmlContent := body.String()
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		return err
	} else if response.StatusCode >= 400 {
		return fmt.Errorf("failed to send email: %s", response.Body)
	}
	return nil
}
