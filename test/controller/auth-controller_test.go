package controller_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/controller"
	"github.com/kuritaeiji/todo-gin-back/mock_service"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type AuthControllerTestSuite struct {
	suite.Suite
	controller      controller.AuthController
	authServiceMock *mock_service.MockAuthService
	rec             *httptest.ResponseRecorder
	ctx             *gin.Context
}

func (suite *AuthControllerTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func (suite *AuthControllerTestSuite) SetupTest() {
	suite.authServiceMock = mock_service.NewMockAuthService(gomock.NewController(suite.T()))
	suite.controller = controller.TestNewAuthController(suite.authServiceMock)
	suite.rec = httptest.NewRecorder()
	suite.ctx, _ = gin.CreateTestContext(suite.rec)
}

func TestAuthController(t *testing.T) {
	suite.Run(t, new(AuthControllerTestSuite))
}

func (suite *AuthControllerTestSuite) TestSuccessLogin() {
	tokenString := "tokenString"
	suite.authServiceMock.EXPECT().Login(suite.ctx).Return(tokenString, nil)
	suite.controller.Login(suite.ctx)

	suite.Equal(200, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), tokenString)
}

func (suite *AuthControllerTestSuite) TestBadLoginWithRecordNotFound() {
	suite.authServiceMock.EXPECT().Login(suite.ctx).Return("", gorm.ErrRecordNotFound)
	suite.controller.Login(suite.ctx)

	suite.Equal(config.RecordNotFoundErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.RecordNotFoundErrorResponse.Json["content"])
}

func (suite *AuthControllerTestSuite) TestBadLoginWithPasswordAuthenticationError() {
	suite.authServiceMock.EXPECT().Login(suite.ctx).Return("", config.PasswordAuthenticationError)
	suite.controller.Login(suite.ctx)

	suite.Equal(config.PasswordAuthenticationErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.PasswordAuthenticationErrorResponse.Json["content"])
}
