package service_test

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang/mock/gomock"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/dto"
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
	var user model.User
	dtoUser := dto.User{Email: email, Password: password}
	dtoUser.Transfer(&user)
	suite.userRepositoryMock.EXPECT().FindByEmail(email).Return(user, nil)
	suite.jwtServiceMock.EXPECT().CreateJWT(user, service.DayFromNowAccessToken).Return(tokenString)

	body := map[string]string{"email": email, "password": password}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/login", strings.NewReader(string(bodyBytes)))
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
	suite.userRepositoryMock.EXPECT().FindByEmail(email).Return(model.User{}, gorm.ErrRecordNotFound)
	body := map[string]string{
		"email": email,
	}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/login", strings.NewReader(string(bodyBytes)))
	req.Header.Add("Content-Type", binding.MIMEJSON)
	suite.ctx.Request = req
	_, err := suite.service.Login(suite.ctx)

	suite.Equal(gorm.ErrRecordNotFound, err)
}

func (suite *AuthServiceTestSuite) TestBadLoginWithPasswordAuthenticationError() {
	user := model.User{Email: email, PasswordDigest: password}
	suite.userRepositoryMock.EXPECT().FindByEmail(email).Return(user, nil)

	body := map[string]string{"email": email, "password": password}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/login", strings.NewReader(string(bodyBytes)))
	req.Header.Add("Content-Type", binding.MIMEJSON)
	suite.ctx.Request = req
	_, err := suite.service.Login(suite.ctx)

	suite.Equal(config.PasswordAuthenticationError, err)
}
