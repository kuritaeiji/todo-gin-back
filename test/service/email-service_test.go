package service_test

import (
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kuritaeiji/todo-gin-back/mock_service"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/service"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stretchr/testify/assert"
)

var (
	emailService    service.EmailService
	jwtServiceMock  *mock_service.MockJWTService
	emailClientMock *mock_service.MockEmailClient
)

func prepareEmailService(t *testing.T) {
	assertion = assert.New(t)
	jwtServiceMock = mock_service.NewMockJWTService(gomock.NewController(t))
	emailClientMock = mock_service.NewMockEmailClient(gomock.NewController(t))
	emailService = service.TestNewEmailService(emailClientMock, jwtServiceMock)
	rec = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(rec)
}

func TestActivationUserEmail(t *testing.T) {
	prepareEmailService(t)
	user := model.User{
		ID:    10,
		Email: "user@example.com",
	}
	token := "token"
	jwtServiceMock.EXPECT().CreateJWT(user.ID, 1).Return(token)
	doFunc := func(msg *mail.SGMailV3) {
		assertion.Equal(os.Getenv("FROM_EMAIL_NAME"), msg.From.Name)
		assertion.Equal(os.Getenv("FROM_EMAIL_ADDRESS"), msg.From.Address)
		assertion.Equal("アカウント有効化リンク", msg.Subject)
		assertion.Equal(user.Email, msg.Personalizations[0].To[0].Address)
		assertion.Contains(msg.Content[0].Value, fmt.Sprintf(`<a href="%v/activate?token=%v">`, os.Getenv("FRONT_ORIGIN"), token))
	}
	emailClientMock.EXPECT().Send(gomock.Any()).Return(&rest.Response{}, nil).Do(doFunc)

	emailService.ActivationUserEmail(user)
}
