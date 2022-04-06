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
	suite.dto = factory.NewDtoUser(&factory.UserConfig{})
	suite.ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
}

func TestUserDto(t *testing.T) {
	suite.Run(t, new(UserDtoTestSuite))
}

func (suite *UserDtoTestSuite) TestSuccessValidation() {
	req := httptest.NewRequest("POST", "/users", factory.CreateUserRequestBody(&factory.UserConfig{}))
	suite.ctx.Request = req
	err := suite.ctx.ShouldBindJSON(&suite.dto)

	verr, _ := err.(validator.ValidationErrors)
	println(verr.Error())

	suite.Nil(err)
}

func (suite *UserDtoTestSuite) TestBadRequiredEmailValidation() {
	req := httptest.NewRequest("POST", "/users", factory.CreateUserRequestBody(&factory.UserConfig{NotUseDefaultValue: true, Email: ""}))
	suite.ctx.Request = req
	err := suite.ctx.ShouldBindJSON(&suite.dto)

	suite.IsType(validator.ValidationErrors{}, err)
	suite.Contains(err.Error(), "required")
}

func (suite *UserDtoTestSuite) TestBadMax100EmailValidation() {
	body := factory.CreateUserRequestBody(&factory.UserConfig{Email: "email@" + strings.Repeat("a", 91) + ".com"})
	req := httptest.NewRequest("POST", "/users", body)
	suite.ctx.Request = req
	err := suite.ctx.ShouldBindJSON(&suite.dto)

	suite.IsType(validator.ValidationErrors{}, err)
	suite.Contains(err.Error(), "max")
}

func (suite *UserDtoTestSuite) TestBadEmailValidation() {
	body := factory.CreateUserRequestBody(&factory.UserConfig{Email: "email"})
	req := httptest.NewRequest("POST", "/users", body)
	suite.ctx.Request = req
	err := suite.ctx.ShouldBindJSON(&suite.dto)

	suite.IsType(validator.ValidationErrors{}, err)
	suite.Contains(err.Error(), "email")
}

func (suite *UserDtoTestSuite) TestBadMin8PasswordValidation() {
	body := factory.CreateUserRequestBody(&factory.UserConfig{Password: "Pas101"})
	req := httptest.NewRequest("POST", "/users", body)
	suite.ctx.Request = req
	err := suite.ctx.ShouldBindJSON(&suite.dto)

	suite.IsType(validator.ValidationErrors{}, err)
	suite.Contains(err.Error(), "min")
}

func (suite *UserDtoTestSuite) TestBadMax50PasswordValidation() {
	body := factory.CreateUserRequestBody(&factory.UserConfig{Password: strings.Repeat("a", 51)})
	req := httptest.NewRequest("POST", "/users", body)
	suite.ctx.Request = req
	err := suite.ctx.ShouldBindJSON(&suite.dto)

	suite.IsType(validator.ValidationErrors{}, err)
	suite.Contains(err.Error(), "max")
}
func (suite *UserDtoTestSuite) TestBadPasswordValidationWithNotHaveUppercase() {
	body := factory.CreateUserRequestBody(&factory.UserConfig{Password: "password1010"})
	req := httptest.NewRequest("POST", "/users", body)
	suite.ctx.Request = req
	err := suite.ctx.ShouldBindJSON(&suite.dto)

	suite.IsType(validator.ValidationErrors{}, err)
	suite.Contains(err.Error(), "password")
}
func (suite *UserDtoTestSuite) TestBadPasswordValidationWithNotHaveUndercase() {
	body := factory.CreateUserRequestBody(&factory.UserConfig{Password: "PASSWORD1010"})
	req := httptest.NewRequest("POST", "/users", body)
	suite.ctx.Request = req
	err := suite.ctx.ShouldBindJSON(&suite.dto)

	suite.IsType(validator.ValidationErrors{}, err)
	suite.Contains(err.Error(), "password")
}

func (suite *UserDtoTestSuite) TestBadPasswordValidationWithNotHaveNumber() {
	body := factory.CreateUserRequestBody(&factory.UserConfig{Password: "Password"})
	req := httptest.NewRequest("POST", "/users", body)
	suite.ctx.Request = req
	err := suite.ctx.ShouldBindJSON(&suite.dto)

	suite.IsType(validator.ValidationErrors{}, err)
	suite.Contains(err.Error(), "password")
}

func (suite *UserDtoTestSuite) TestTransfer() {
	var user model.User
	suite.dto.Transfer(&user)

	suite.Contains(user.Email, factory.DefaultEmail)
	suite.True(user.Authenticate(factory.DefualtPassword))
}
