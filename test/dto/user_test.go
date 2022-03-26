package dto_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kuritaeiji/todo-gin-back/dto"
	"github.com/kuritaeiji/todo-gin-back/factory"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/validators"
	"github.com/stretchr/testify/suite"
)

type UserDtoTestSuite struct {
	suite.Suite
	dto dto.User
	ctx *gin.Context
}

func (suite *UserDtoTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	validators.Init()
}

func (suite *UserDtoTestSuite) SetupTest() {
	suite.dto = factory.NewDtoUser(factory.UserConfig{})
	suite.ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
}

func TestUserDto(t *testing.T) {
	suite.Run(t, new(UserDtoTestSuite))
}

func (suite *UserDtoTestSuite) TestSuccessValidation() {
	req := httptest.NewRequest("POST", "/users", factory.CreateUserRequestBody(factory.DefaultEmail, factory.DefualtPassword))
	suite.ctx.Request = req
	err := suite.ctx.ShouldBindJSON(&suite.dto)

	suite.Nil(err)
}

func (suite *UserDtoTestSuite) TestBadRequiredEmailValidation() {
	req := httptest.NewRequest("POST", "/users", factory.CreateUserRequestBody("", factory.DefualtPassword))
	suite.ctx.Request = req
	err := suite.ctx.ShouldBindJSON(&suite.dto)

	suite.IsType(validator.ValidationErrors{}, err)
	suite.Contains(err.Error(), "required")
}

func (suite *UserDtoTestSuite) TestBadMax100EmailValidation() {
	body := factory.CreateUserRequestBody("email@"+strings.Repeat("a", 91)+".com", factory.DefualtPassword)
	req := httptest.NewRequest("POST", "/users", body)
	suite.ctx.Request = req
	err := suite.ctx.ShouldBindJSON(&suite.dto)

	suite.IsType(validator.ValidationErrors{}, err)
	suite.Contains(err.Error(), "max")
}

func (suite *UserDtoTestSuite) TestBadEmailValidation() {
	body := factory.CreateUserRequestBody("email", factory.DefualtPassword)
	req := httptest.NewRequest("POST", "/users", body)
	suite.ctx.Request = req
	err := suite.ctx.ShouldBindJSON(&suite.dto)

	suite.IsType(validator.ValidationErrors{}, err)
	suite.Contains(err.Error(), "email")
}

func (suite *UserDtoTestSuite) TestBadMin8PasswordValidation() {
	body := factory.CreateUserRequestBody(factory.DefualtPassword, "Pass101")
	req := httptest.NewRequest("POST", "/users", body)
	suite.ctx.Request = req
	err := suite.ctx.ShouldBindJSON(&suite.dto)

	suite.IsType(validator.ValidationErrors{}, err)
	suite.Contains(err.Error(), "min")
}

func (suite *UserDtoTestSuite) TestBadMax50PasswordValidation() {
	body := factory.CreateUserRequestBody(factory.DefualtPassword, strings.Repeat("a", 51))
	req := httptest.NewRequest("POST", "/users", body)
	suite.ctx.Request = req
	err := suite.ctx.ShouldBindJSON(&suite.dto)

	suite.IsType(validator.ValidationErrors{}, err)
	suite.Contains(err.Error(), "max")
}
func (suite *UserDtoTestSuite) TestBadPasswordValidationWithNotHaveUppercase() {
	body := factory.CreateUserRequestBody(factory.DefualtPassword, "password1010")
	req := httptest.NewRequest("POST", "/users", body)
	suite.ctx.Request = req
	err := suite.ctx.ShouldBindJSON(&suite.dto)

	suite.IsType(validator.ValidationErrors{}, err)
	suite.Contains(err.Error(), "password")
}
func (suite *UserDtoTestSuite) TestBadPasswordValidationWithNotHaveUndercase() {
	body := factory.CreateUserRequestBody(factory.DefualtPassword, "PASSWORD1010")
	req := httptest.NewRequest("POST", "/users", body)
	suite.ctx.Request = req
	err := suite.ctx.ShouldBindJSON(&suite.dto)

	suite.IsType(validator.ValidationErrors{}, err)
	suite.Contains(err.Error(), "password")
}

func (suite *UserDtoTestSuite) TestBadPasswordValidationWithNotHaveNumber() {
	body := factory.CreateUserRequestBody(factory.DefualtPassword, "Password")
	req := httptest.NewRequest("POST", "/users", body)
	suite.ctx.Request = req
	err := suite.ctx.ShouldBindJSON(&suite.dto)

	suite.IsType(validator.ValidationErrors{}, err)
	suite.Contains(err.Error(), "password")
}

func (suite *UserDtoTestSuite) TestTransfer() {
	var user model.User
	suite.dto.Transfer(&user)

	suite.Equal(factory.DefaultEmail, user.Email)
	suite.True(user.Authenticate(factory.DefualtPassword))
}
