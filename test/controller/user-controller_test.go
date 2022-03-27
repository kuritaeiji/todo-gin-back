package controller_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/golang/mock/gomock"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/controller"
	"github.com/kuritaeiji/todo-gin-back/mock_service"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UserControllerTestSuite struct {
	suite.Suite
	controller       controller.UserController
	userServiceMock  *mock_service.MockUserService
	emailServiceMock *mock_service.MockEmailService
	ctx              *gin.Context
	rec              *httptest.ResponseRecorder
}

func (suite *UserControllerTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func (suite *UserControllerTestSuite) SetupTest() {
	suite.userServiceMock = mock_service.NewMockUserService(gomock.NewController(suite.T()))
	suite.emailServiceMock = mock_service.NewMockEmailService(gomock.NewController(suite.T()))
	suite.controller = controller.TestNewUserController(suite.userServiceMock, suite.emailServiceMock)
	suite.rec = httptest.NewRecorder()
	suite.ctx, _ = gin.CreateTestContext(suite.rec)
}

func TestUesrController(t *testing.T) {
	suite.Run(t, &UserControllerTestSuite{})
}

func (suite *UserControllerTestSuite) TestSuccessCreate() {
	var user model.User
	suite.userServiceMock.EXPECT().Create(suite.ctx).Return(user, nil)
	suite.emailServiceMock.EXPECT().ActivationUserEmail(user)
	suite.controller.Create(suite.ctx)

	suite.Equal(200, suite.rec.Code)
}

func (suite *UserControllerTestSuite) TestBadCreateWithValidationError() {
	var verr validator.ValidationErrors
	var err error = verr
	suite.userServiceMock.EXPECT().Create(suite.ctx).Return(model.User{}, err)

	suite.controller.Create(suite.ctx)
	suite.Equal(config.ValidationErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.ValidationErrorResponse.Json["content"])
}

func (suite *UserControllerTestSuite) TestBadCreateWithNotUniqueUser() {
	err := config.UniqueUserError
	suite.userServiceMock.EXPECT().Create(suite.ctx).Return(model.User{}, err)

	suite.controller.Create(suite.ctx)
	suite.Equal(config.UniqueUserErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.UniqueUserErrorResponse.Json["content"])
}

func (suite *UserControllerTestSuite) TestBadCreateWithnEmailClientError() {
	err := config.EmailClientError
	var user model.User
	suite.userServiceMock.EXPECT().Create(suite.ctx).Return(user, nil)
	suite.emailServiceMock.EXPECT().ActivationUserEmail(user).Return(err)
	suite.controller.Create(suite.ctx)

	suite.Equal(config.EmailClientErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.EmailClientErrorResponse.Json["content"])
}

func (suite *UserControllerTestSuite) TestTrueIsUnique() {
	suite.userServiceMock.EXPECT().IsUnique(suite.ctx).Return(true, nil)
	suite.controller.IsUnique(suite.ctx)
	suite.Equal(200, suite.rec.Code)
}

func (suite *UserControllerTestSuite) TestFalseIsUnique() {
	suite.userServiceMock.EXPECT().IsUnique(suite.ctx).Return(false, nil)
	suite.controller.IsUnique(suite.ctx)
	suite.Equal(config.UniqueUserErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.UniqueUserErrorResponse.Json["content"])
}

func (suite *UserControllerTestSuite) TestSuccessActivate() {
	suite.userServiceMock.EXPECT().Activate(suite.ctx).Return(nil)
	suite.controller.Activate(suite.ctx)
	suite.Equal(200, suite.rec.Code)
}

func (suite *UserControllerTestSuite) TestBadActivateWithJWTExpired() {
	jwtErr := jwt.NewValidationError("", jwt.ValidationErrorExpired)
	var err error = jwtErr
	suite.userServiceMock.EXPECT().Activate(suite.ctx).Return(err)
	suite.controller.Activate(suite.ctx)
	suite.Equal(config.JWTExpiredErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.JWTExpiredErrorResponse.Json["content"])
}

func (suite *UserControllerTestSuite) TestBadActivateWithJWTValidationError() {
	jwtErr := jwt.NewValidationError("", jwt.ValidationErrorClaimsInvalid)
	var err error = jwtErr
	suite.userServiceMock.EXPECT().Activate(suite.ctx).Return(err)
	suite.controller.Activate(suite.ctx)
	suite.Equal(config.JWTValidationErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.JWTValidationErrorResponse.Json["content"])
}

func (suite *UserControllerTestSuite) TestBadActivateWithRecordNotFound() {
	suite.userServiceMock.EXPECT().Activate(suite.ctx).Return(gorm.ErrRecordNotFound)
	suite.controller.Activate(suite.ctx)
	suite.Equal(config.RecordNotFoundErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.RecordNotFoundErrorResponse.Json["content"])
}

func (suite *UserControllerTestSuite) TestBadActivateWithAlreadyActivatedUser() {
	suite.userServiceMock.EXPECT().Activate(suite.ctx).Return(config.AlreadyActivatedUserError)
	suite.controller.Activate(suite.ctx)
	suite.Equal(config.AlreadyActivatedUserErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.AlreadyActivatedUserErrorResponse.Json["content"])
}
