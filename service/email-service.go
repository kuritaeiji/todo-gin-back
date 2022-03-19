package service

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailService interface {
	ActivationUserEmail(model.User)
}

type emailService struct {
	client     *sendgrid.Client
	jwtService JWTService
	from       string
}

func NewEmailService() EmailService {
	return &emailService{
		sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY")),
		NewJWTService(),
		os.Getenv("EMAIL"),
	}
}

func (s *emailService) ActivationUserEmail(user model.User) {
	from := mail.NewEmail("", s.from)
	subject := "アカウント有効化リンク"
	to := mail.NewEmail("", user.Email)
	message := mail.NewSingleEmail(from, subject, to, "", s.activationHTML(user))
	_, err := s.client.Send(message)
	if err != nil {
		gin.DefaultWriter.Write([]byte(fmt.Sprintf("Failed to send activation user email\n%v", err.Error())))
	}
}

func (s *emailService) activationHTML(user model.User) string {
	token := s.jwtService.CreateJWT(user.ID, 1)
	html := template.Must(template.ParseFiles("template/activation-user.html"))
	pr, pw := io.Pipe()
	go func() {
		html.Execute(pw, fmt.Sprintf("%v/activate?token=%v", os.Getenv("FRONT_ORIGIN"), token))
		pw.Close()
	}()
	byteSlice, _ := ioutil.ReadAll(pr)
	return string(byteSlice)
}
