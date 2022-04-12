package dto_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kuritaeiji/todo-gin-back/dto"
	"github.com/kuritaeiji/todo-gin-back/factory"
	"github.com/kuritaeiji/todo-gin-back/validators"
	"github.com/stretchr/testify/suite"
)

type ListDtoTestSuite struct {
	suite.Suite
	dto dto.List
	ctx *gin.Context
}

func (suite *ListDtoTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	validators.Init()
}

func (suite *ListDtoTestSuite) SetupTest() {
	suite.dto = factory.NewDtoList(&factory.ListConfig{})
	suite.ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
}

func TestListDto(t *testing.T) {
	suite.Run(t, new(ListDtoTestSuite))
}

func (suite *ListDtoTestSuite) TestSuccessValidation() {
	req := httptest.NewRequest("POST", "/api/lists", factory.CreateListRequestBody(&factory.ListConfig{}))
	suite.ctx.Request = req
	err := suite.ctx.ShouldBindJSON(&suite.dto)

	suite.Nil(err)
}

func (suite *ListDtoTestSuite) TestBadValidationWithTitleRequired() {
	body := factory.CreateListRequestBody(&factory.ListConfig{NotUseDefaultValue: true})
	req := httptest.NewRequest("POST", "/api/lists", body)
	suite.ctx.Request = req
	err := suite.ctx.ShouldBindJSON(&suite.dto)
	verr, _ := err.(validator.ValidationErrors)

	suite.Equal("required", verr[0].Tag())
	suite.Equal("Title", verr[0].Field())
}

func (suite *ListDtoTestSuite) TestBadValidationWithTitleMax50() {
	body := factory.CreateListRequestBody(&factory.ListConfig{Title: strings.Repeat("a", 51)})
	req := httptest.NewRequest("POST", "/api/lists", body)
	suite.ctx.Request = req
	err := suite.ctx.ShouldBindJSON(&suite.dto)
	verr, _ := err.(validator.ValidationErrors)

	suite.Equal("max", verr[0].Tag())
	suite.Equal("Title", verr[0].Field())
}

func (suite *ListDtoTestSuite) TestBadValidationWithIndexGreaterThenEqual0() {
	body := factory.CreateListRequestBody(&factory.ListConfig{Index: -1})
	req := httptest.NewRequest("POST", "/api/lists", body)
	suite.ctx.Request = req
	err := suite.ctx.ShouldBindJSON(&suite.dto)
	verr, _ := err.(validator.ValidationErrors)

	suite.Equal("gte", verr[0].Tag())
	suite.Equal("Index", verr[0].Field())
}
