package middleware_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/golang/mock/gomock"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/middleware"
	"github.com/kuritaeiji/todo-gin-back/mock_repository"
	"github.com/kuritaeiji/todo-gin-back/mock_service"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/service"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type AuthMiddlewareTestSuite struct {
	suite.Suite
	middleware         middleware.AuthMiddleware
	jwtServiceMock     *mock_service.MockJWTService
	userRepositoryMock *mock_repository.MockUserRepository
	rec                *httptest.ResponseRecorder
	ctx                *gin.Context
}

func (suite *AuthMiddlewareTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func (suite *AuthMiddlewareTestSuite) SetupTest() {
	suite.jwtServiceMock = mock_service.NewMockJWTService(gomock.NewController(suite.T()))
	suite.userRepositoryMock = mock_repository.NewMockUserRepository(gomock.NewController(suite.T()))
	suite.middleware = middleware.TestNewAuthMiddleware(suite.jwtServiceMock, suite.userRepositoryMock)
	suite.rec = httptest.NewRecorder()
	suite.ctx, _ = gin.CreateTestContext(suite.rec)
}

func TestAuthMiddleware(t *testing.T) {
	suite.Run(t, new(AuthMiddlewareTestSuite))
}

func (suite *AuthMiddlewareTestSuite) TestSuccessAuth() {
	var user model.User
	accessToken := "token"
	suite.jwtServiceMock.EXPECT().VerifyJWT(accessToken).Return(
		&service.UserClaim{ID: user.ID, StandardClaims: jwt.StandardClaims{}},
		nil,
	)
	suite.userRepositoryMock.EXPECT().Find(user.ID).Return(user, nil)
	req := httptest.NewRequest("POST", "/users", nil)
	req.Header.Add(config.TokenHeader, config.Bearer+accessToken)
	suite.ctx.Request = req
	suite.middleware.Auth(suite.ctx)

	currentUser := suite.ctx.MustGet("currentUser").(model.User)
	suite.Equal(user, currentUser)
}

func (suite *AuthMiddlewareTestSuite) TestBadAuthWithExpiredJWT() {
	accessToken := "token"
	verr := jwt.NewValidationError("expired", jwt.ValidationErrorExpired)
	var err error = verr
	suite.jwtServiceMock.EXPECT().VerifyJWT(accessToken).Return(&service.UserClaim{}, err)
	req := httptest.NewRequest("POST", "/users", nil)
	req.Header.Add(config.TokenHeader, config.Bearer+accessToken)
	suite.ctx.Request = req
	suite.middleware.Auth(suite.ctx)

	suite.Equal(config.JWTExpiredErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.JWTExpiredErrorResponse.Json["content"])
}

func (suite *AuthMiddlewareTestSuite) TestBadAuthWithJWTValidationError() {
	accessToken := "token"
	verr := jwt.NewValidationError("expired", jwt.ValidationErrorAudience)
	var err error = verr
	suite.jwtServiceMock.EXPECT().VerifyJWT(accessToken).Return(&service.UserClaim{}, err)
	req := httptest.NewRequest("POST", "/users", nil)
	req.Header.Add(config.TokenHeader, config.Bearer+accessToken)
	suite.ctx.Request = req
	suite.middleware.Auth(suite.ctx)

	suite.Equal(config.NotLoggedInErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.NotLoggedInErrorResponse.Json["content"])
}

func (suite *AuthMiddlewareTestSuite) TestBadAuthWithNotRecordFound() {
	accessToken := "token"
	suite.jwtServiceMock.EXPECT().VerifyJWT(accessToken).Return(&service.UserClaim{}, nil)
	suite.userRepositoryMock.EXPECT().Find(0).Return(model.User{}, gorm.ErrRecordNotFound)
	req := httptest.NewRequest("POST", "/users", nil)
	req.Header.Add(config.TokenHeader, config.Bearer+accessToken)
	suite.ctx.Request = req
	suite.middleware.Auth(suite.ctx)

	suite.Equal(config.NotLoggedInErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.NotLoggedInErrorResponse.Json["content"])
}

func (suite *AuthMiddlewareTestSuite) TestSuccessGuest() {
	req := httptest.NewRequest("POST", "/users", nil)
	suite.ctx.Request = req
	suite.middleware.Guest(suite.ctx)

	suite.Equal(200, suite.rec.Code)
}

func (suite *AuthMiddlewareTestSuite) TestBadGuestWithLoggedIn() {
	req := httptest.NewRequest("POST", "/users", nil)
	req.Header.Add(config.TokenHeader, config.Bearer+"token")
	suite.ctx.Request = req
	suite.middleware.Guest(suite.ctx)

	suite.Equal(config.GuestErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.GuestErrorResponse.Json["content"])
}
