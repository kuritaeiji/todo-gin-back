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
	"github.com/kuritaeiji/todo-gin-back/gateway"
	"github.com/kuritaeiji/todo-gin-back/model"
)

type EmailService interface {
	ActivationUserEmail(model.User) error
}

type emailService struct {
	gateway    gateway.EmailGateway
	jwtService JWTService
}

func NewEmailService() EmailService {
	return &emailService{
		gateway.NewEmailGateway(),
		NewJWTService(),
	}
}

func (s *emailService) ActivationUserEmail(user model.User) error {
	err := s.gateway.Send(user.Email, "アカウント有効化リンク", s.activationHTML(user))
	if err != nil {
		gin.DefaultWriter.Write([]byte(fmt.Sprintf("Failed to send activation user email\n%v", err.Error())))
		return config.EmailClientError
	}
	return nil
}

func (s *emailService) activationHTML(user model.User) string {
	token := s.jwtService.CreateJWT(user, DayFromNowActivateUserToken)
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
func TestNewEmailService(gateway gateway.EmailGateway, jwtService JWTService) EmailService {
	return &emailService{
		gateway:    gateway,
		jwtService: jwtService,
	}
}
