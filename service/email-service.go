package service

// mockgen -source=service/email-service.go -destination=./mock_service/email-service.go

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailService interface {
	ActivationUserEmail(model.User)
}

type EmailClient interface {
	Send(*mail.SGMailV3) (*rest.Response, error)
}

type emailService struct {
	client     EmailClient
	jwtService JWTService
	from       *mail.Email
}

func NewEmailService() EmailService {
	return &emailService{
		sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY")),
		NewJWTService(),
		mail.NewEmail(os.Getenv("FROM_EMAIL_NAME"), os.Getenv("FROM_EMAIL_ADDRESS")),
	}
}

func (s *emailService) ActivationUserEmail(user model.User) {
	subject := "アカウント有効化リンク"
	to := mail.NewEmail("", user.Email)
	message := mail.NewSingleEmail(s.from, subject, to, "", s.activationHTML(user))
	_, err := s.client.Send(message)
	if err != nil {
		gin.DefaultWriter.Write([]byte(fmt.Sprintf("Failed to send activation user email\n%v", err.Error())))
	}
}

func (s *emailService) activationHTML(user model.User) string {
	token := s.jwtService.CreateJWT(user.ID, 1)
	html := template.Must(template.ParseFiles(fmt.Sprintf("%v/template/activation-user.html", config.WorkDir)))
	pr, pw := io.Pipe()
	go func() {
		html.Execute(pw, fmt.Sprintf("%v/activate?token=%v", os.Getenv("FRONT_ORIGIN"), token))
		pw.Close()
	}()
	byteSlice, _ := ioutil.ReadAll(pr)
	return string(byteSlice)
}

// test用
func TestNewEmailService(client EmailClient, jwtService JWTService) EmailService {
	return &emailService{
		client:     client,
		jwtService: jwtService,
		from:       mail.NewEmail(os.Getenv("FROM_EMAIL_NAME"), os.Getenv("FROM_EMAIL_ADDRESS")),
	}
}
