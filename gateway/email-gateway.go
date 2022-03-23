package gateway

// mockgen -source=gateway/email-gateway.go -destination=mock_gateway/email-gateway.go

import (
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailGateway interface {
	Send(to, subject, htmlString string) error
}

type emailGateway struct {
	client   *sendgrid.Client
	fromMail *mail.Email
}

func NewEmailGateway() EmailGateway {
	return &emailGateway{
		client:   sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY")),
		fromMail: mail.NewEmail(os.Getenv("FROM_EMAIL_NAME"), os.Getenv("FROM_EMAIL_ADDRESS")),
	}
}

func (gateway *emailGateway) Send(to, subject, htmlString string) error {
	toMail := mail.NewEmail("", to)
	message := mail.NewSingleEmail(gateway.fromMail, subject, toMail, "", htmlString)
	_, err := gateway.client.Send(message)
	return err
}
