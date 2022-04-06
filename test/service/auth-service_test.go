package service_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang/mock/gomock"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/factory"
	"github.com/kuritaeiji/todo-gin-back/mock_repository"
	"github.com/kuritaeiji/todo-gin-back/mock_service"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/service"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type AuthServiceTestSuite struct {
	suite.Suite
	service            service.AuthService
	userRepositoryMock *mock_repository.MockUserRepository
	jwtServiceMock     *mock_service.MockJWTService
	rec                *httptest.ResponseRecorder
	ctx                *gin.Context
}

func (suite *AuthServiceTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func (suite *AuthServiceTestSuite) SetupTest() {
	suite.userRepositoryMock = mock_repository.NewMockUserRepository(gomock.NewController(suite.T()))
	suite.jwtServiceMock = mock_service.NewMockJWTService(gomock.NewController(suite.T()))
	suite.service = service.TestNewAuthService(suite.userRepositoryMock, suite.jwtServiceMock)
	suite.rec = httptest.NewRecorder()
	suite.ctx, _ = gin.CreateTestContext(suite.rec)
}

func TestAuthService(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}

func (suite *AuthServiceTestSuite) TestSuccessLogin() {
	user := factory.NewUser(factory.UserConfig{})
	tokenString := factory.CreateAccessToken(user)
	suite.userRepositoryMock.EXPECT().FindByEmail(user.Email).Return(user, nil)
	suite.jwtServiceMock.EXPECT().CreateJWT(user, service.DayFromNowAccessToken).Return(tokenString)

	req := httptest.NewRequest("POST", "/login", factory.CreateUserRequestBody(factory.UserConfig{Email: user.Email}))
	req.Header.Add("Content-Type", binding.MIMEJSON)
	suite.ctx.Request = req
	token, err := suite.service.Login(suite.ctx)

	suite.Equal(tokenString, token)
	suite.Nil(err)
}

func (suite *AuthServiceTestSuite) TestBadLoginWithCannotBind() {
	req := httptest.NewRequest("POST", "/login", nil)
	req.Header.Add("Content-Type", binding.MIMEJSON)
	suite.ctx.Request = req
	_, err := suite.service.Login(suite.ctx)

	suite.Error(err)
}

func (suite *AuthServiceTestSuite) TestBadLoginWithRecordNotFound() {
	req := httptest.NewRequest("POST", "/login", factory.CreateUserRequestBody(factory.UserConfig{}))
	req.Header.Add("Content-Type", binding.MIMEJSON)
	suite.ctx.Request = req
	suite.userRepositoryMock.EXPECT().FindByEmail(gomock.Any()).Return(model.User{}, gorm.ErrRecordNotFound).Do(func(email string) {
		suite.Contains(email, factory.DefaultEmail)
	})
	_, err := suite.service.Login(suite.ctx)

	suite.Equal(gorm.ErrRecordNotFound, err)
}

func (suite *AuthServiceTestSuite) TestBadLoginWithPasswordAuthenticationError() {
	user := factory.NewUser(factory.UserConfig{})
	suite.userRepositoryMock.EXPECT().FindByEmail(user.Email).Return(user, nil)

	req := httptest.NewRequest("POST", "/login", factory.CreateUserRequestBody(factory.UserConfig{Email: user.Email, Password: "invalid password"}))
	req.Header.Add("Content-Type", binding.MIMEJSON)
	suite.ctx.Request = req
	_, err := suite.service.Login(suite.ctx)

	suite.Equal(config.PasswordAuthenticationError, err)
}
