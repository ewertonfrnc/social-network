package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridmailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendGrid(fromEmail, apiKey string) *SendGridmailer {
	client := sendgrid.NewSendClient(apiKey)

	return &SendGridmailer{
		fromEmail,
		apiKey,
		client,
	}
}

func (m *SendGridmailer) Send(templateFile, toUsername, toEmail string, data any, isSandbox bool) error {
	from := mail.NewEmail(fromEmail, m.fromEmail)
	to := mail.NewEmail(toUsername, toEmail)

	// Parsing template
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return err
	}

	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())

	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})

	for i := range maxRetries {
		response, err := m.client.Send(message)
		if err != nil {
			log.Printf("Failed to send email to %v. Attempt %d of %d", toEmail, i+1, maxRetries)
			log.Printf("[ERROR]: %v", err)
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		if response.StatusCode >= 400 {
			log.Printf("Failed to send email to %v, status %d. Attempt %d of %d", toEmail, response.StatusCode, i+1, maxRetries)
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		log.Printf("Email sent with status code %v", response.StatusCode)
		return nil
	}

	return fmt.Errorf("Failed to send email after %d retries", maxRetries)
}
