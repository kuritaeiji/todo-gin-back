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
	var userConfig factory.UserConfig
	user := factory.NewUser(&userConfig)
	tokenString := factory.CreateAccessToken(user)
	suite.userRepositoryMock.EXPECT().FindByEmail(userConfig.Email).Return(user, nil)
	suite.jwtServiceMock.EXPECT().CreateJWT(user, service.DayFromNowAccessToken).Return(tokenString)

	req := httptest.NewRequest("POST", "/login", factory.CreateUserRequestBody(&userConfig))
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
	var userConfig factory.UserConfig
	req := httptest.NewRequest("POST", "/login", factory.CreateUserRequestBody(&userConfig))
	req.Header.Add("Content-Type", binding.MIMEJSON)
	suite.ctx.Request = req
	suite.userRepositoryMock.EXPECT().FindByEmail(userConfig.Email).Return(model.User{}, gorm.ErrRecordNotFound)
	_, err := suite.service.Login(suite.ctx)

	suite.Equal(gorm.ErrRecordNotFound, err)
}

func (suite *AuthServiceTestSuite) TestBadLoginWithPasswordAuthenticationError() {
	var userConfig factory.UserConfig
	user := factory.NewUser(&userConfig)
	suite.userRepositoryMock.EXPECT().FindByEmail(user.Email).Return(user, nil)

	userConfig.Password = "invalid password"
	req := httptest.NewRequest("POST", "/login", factory.CreateUserRequestBody(&userConfig))
	req.Header.Add("Content-Type", binding.MIMEJSON)
	suite.ctx.Request = req
	_, err := suite.service.Login(suite.ctx)

	suite.Equal(config.PasswordAuthenticationError, err)
}
