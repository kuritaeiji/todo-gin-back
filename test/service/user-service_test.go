package service_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/golang/mock/gomock"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/mock_repository"
	"github.com/kuritaeiji/todo-gin-back/mock_service"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/service"
	"github.com/kuritaeiji/todo-gin-back/validators"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceTestSuite struct {
	suite.Suite
	service            service.UserService
	userRepositoryMock *mock_repository.MockUserRepository
	jwtServiceMock     *mock_service.MockJWTService
	rec                *httptest.ResponseRecorder
	ctx                *gin.Context
}

func (suite *UserServiceTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	validators.Init()
}

func (suite *UserServiceTestSuite) SetupTest() {
	suite.userRepositoryMock = mock_repository.NewMockUserRepository(gomock.NewController(suite.T()))
	suite.jwtServiceMock = mock_service.NewMockJWTService(gomock.NewController(suite.T()))
	suite.service = service.TestNewUserService(suite.jwtServiceMock, suite.userRepositoryMock)
	suite.rec = httptest.NewRecorder()
	suite.ctx, _ = gin.CreateTestContext(suite.rec)
}

func TestUserServiceTest(t *testing.T) {
	suite.Run(t, &UserServiceTestSuite{})
}

func (suite *UserServiceTestSuite) TestSuccessCreate() {
	body := map[string]string{
		"email":    email,
		"password": password,
	}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/users", strings.NewReader(string(bodyBytes)))
	req.Header.Add("Content-Type", binding.MIMEJSON)
	suite.ctx.Request = req
	suite.userRepositoryMock.EXPECT().Create(gomock.Any()).Return(nil).Do(func(user *model.User) {
		suite.Equal(email, user.Email)
		suite.Nil(bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password)))
	})
	user, err := suite.service.Create(suite.ctx)

	suite.Nil(err)
	suite.Equal(email, user.Email)
	suite.Nil(bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password)))
}

func (suite *UserServiceTestSuite) TestBadCreateWithValidation() {
	req := httptest.NewRequest("POST", "/users", strings.NewReader(`{"email":"", "password":""}`))
	req.Header.Add("Content-Type", binding.MIMEJSON)
	suite.ctx.Request = req
	_, err := suite.service.Create(suite.ctx)
	suite.IsType(validator.ValidationErrors{}, err)
}

func (suite *UserServiceTestSuite) TestBadCreateWithDBError() {
	body := map[string]string{
		"email":    email,
		"password": password,
	}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/users", strings.NewReader(string(bodyBytes)))
	req.Header.Add("Content-Type", binding.MIMEJSON)
	suite.ctx.Request = req
	err := errors.New("error")
	suite.userRepositoryMock.EXPECT().Create(gomock.Any()).Return(err).Do(func(user *model.User) {
		suite.Equal(email, user.Email)
		suite.Nil(bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password)))
	})
	_, rerr := suite.service.Create(suite.ctx)

	suite.Equal(err, rerr)
}

func (suite *UserServiceTestSuite) TestTrueIsUnique() {
	req := httptest.NewRequest("GET", fmt.Sprintf("/users/unique?email=%v", email), nil)
	suite.ctx.Request = req
	suite.userRepositoryMock.EXPECT().IsUnique(email).Return(true, nil)
	result, err := suite.service.IsUnique(suite.ctx)

	suite.True(result)
	suite.Nil(err)
}

func (suite *UserServiceTestSuite) TestFalseIsUnique() {
	req := httptest.NewRequest("GET", fmt.Sprintf("/users/unique?email=%v", email), nil)
	suite.ctx.Request = req
	suite.userRepositoryMock.EXPECT().IsUnique(email).Return(false, nil)
	result, err := suite.service.IsUnique(suite.ctx)

	suite.False(result)
	suite.Nil(err)
}

func (suite *UserServiceTestSuite) TestBadIsUniqueWithDBError() {
	req := httptest.NewRequest("GET", fmt.Sprintf("/users/unique?email=%v", email), nil)
	suite.ctx.Request = req
	suite.userRepositoryMock.EXPECT().IsUnique(email).Return(false, errors.New("db error"))
	result, err := suite.service.IsUnique(suite.ctx)

	suite.False(result)
	suite.Error(err)
}

func (suite *UserServiceTestSuite) TestSuccessActivate() {
	user := model.User{}
	suite.jwtServiceMock.EXPECT().VerifyJWT(tokenString).Return(claim, nil)
	suite.userRepositoryMock.EXPECT().Find(id).Return(user, nil)
	suite.userRepositoryMock.EXPECT().Activate(&user).Return(nil)

	req := httptest.NewRequest("PUT", fmt.Sprintf("/users/activate?token=%v", tokenString), nil)
	suite.ctx.Request = req
	err := suite.service.Activate(suite.ctx)

	suite.Nil(err)
}

func (suite *UserServiceTestSuite) TestBadActivateWithJWTValidationError() {
	err := errors.New("jwt validation error")
	suite.jwtServiceMock.EXPECT().VerifyJWT(tokenString).Return(claim, err)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/users/activate?token=%v", tokenString), nil)
	suite.ctx.Request = req
	rerr := suite.service.Activate(suite.ctx)

	suite.Equal(err, rerr)
}

func (suite *UserServiceTestSuite) TestBadActivateWithRecordNotFound() {
	err := errors.New("record not found")
	suite.jwtServiceMock.EXPECT().VerifyJWT(tokenString).Return(claim, nil)
	suite.userRepositoryMock.EXPECT().Find(claim.ID).Return(model.User{}, err)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/users/activate?token=%v", tokenString), nil)
	suite.ctx.Request = req
	rerr := suite.service.Activate(suite.ctx)

	suite.Equal(err, rerr)
}

func (suite *UserServiceTestSuite) TestBadActivateWithAlreadyActivated() {
	suite.jwtServiceMock.EXPECT().VerifyJWT(tokenString).Return(claim, nil)
	suite.userRepositoryMock.EXPECT().Find(claim.ID).Return(model.User{Activated: true}, nil)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/users/activate?token=%v", tokenString), nil)
	suite.ctx.Request = req
	err := suite.service.Activate(suite.ctx)

	suite.Equal(config.AlreadyActivatedUserError, err)
}

func (suite *UserServiceTestSuite) TestBadActivateWithDBError() {
	err := errors.New("db error")
	user := model.User{}
	suite.jwtServiceMock.EXPECT().VerifyJWT(tokenString).Return(claim, nil)
	suite.userRepositoryMock.EXPECT().Find(claim.ID).Return(user, nil)
	suite.userRepositoryMock.EXPECT().Activate(&user).Return(err)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/users/activate?token=%v", tokenString), nil)
	suite.ctx.Request = req
	rerr := suite.service.Activate(suite.ctx)

	suite.Equal(err, rerr)
}
