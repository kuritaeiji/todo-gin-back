package service_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/mock_gateway"
	"github.com/kuritaeiji/todo-gin-back/mock_service"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/service"
	"github.com/stretchr/testify/suite"
)

type EmailServiceTestSuite struct {
	suite.Suite
	service          service.EmailService
	emailGatewayMock *mock_gateway.MockEmailGateway
	jwtServiceMock   *mock_service.MockJWTService
}

func (suite *EmailServiceTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	config.Init()
}

func (suite *EmailServiceTestSuite) SetupTest() {
	suite.emailGatewayMock = mock_gateway.NewMockEmailGateway(gomock.NewController(suite.T()))
	suite.jwtServiceMock = mock_service.NewMockJWTService(gomock.NewController(suite.T()))
	suite.service = service.TestNewEmailService(suite.emailGatewayMock, suite.jwtServiceMock)
}

func TestEmailServiceSuite(t *testing.T) {
	suite.Run(t, &EmailServiceTestSuite{})
}

func (suite *EmailServiceTestSuite) TestSuccessActivationUserEmail() {
	var user model.User
	token := "token"
	doFunc := func(to, subject, htmlString string) {
		suite.Contains(htmlString, fmt.Sprintf(`<a href="%v/activate?token=%v`, os.Getenv("FRONT_ORIGIN"), token))
	}
	suite.emailGatewayMock.EXPECT().Send(user.Email, "アカウント有効化リンク", gomock.Any()).Return(nil).Do(doFunc)
	suite.jwtServiceMock.EXPECT().CreateJWT(user, 1).Return(token)
	err := suite.service.ActivationUserEmail(user)

	suite.Nil(err)
}

func (suite *EmailServiceTestSuite) TestBadActivationUserEmailWithEmailGatewayError() {
	var user model.User
	token := "token"
	err := errors.New("email client error")
	suite.jwtServiceMock.EXPECT().CreateJWT(user, 1).Return(token)
	suite.emailGatewayMock.EXPECT().Send(user.Email, "アカウント有効化リンク", gomock.Any()).Return(err)
	rerr := suite.service.ActivationUserEmail(user)

	suite.Equal(config.EmailClientError, rerr)
}
