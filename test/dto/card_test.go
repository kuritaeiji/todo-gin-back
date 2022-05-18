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

type CardDtoTestSuite struct {
	suite.Suite
	ctx *gin.Context
	dto *dto.Card
}

func (suite *CardDtoTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	validators.Init()
}

func (suite *CardDtoTestSuite) SetupTest() {
	suite.ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
	suite.dto = &dto.Card{}
}

func TestCardDto(t *testing.T) {
	suite.Run(t, new(CardDtoTestSuite))
}

func (suite *CardDtoTestSuite) TestSuccessValidation() {
	cardConfig := &factory.CardConfig{}
	req := httptest.NewRequest("POST", "/", factory.CreateCardRequestBody(cardConfig))
	suite.ctx.Request = req
	err := suite.ctx.ShouldBindJSON(suite.dto)
	suite.Nil(err)
}

func (suite *CardDtoTestSuite) TestBadValidationWithTitleRequired() {
	cardConfig := &factory.CardConfig{NotUseDefaultValue: true}
	req := httptest.NewRequest("POST", "/", factory.CreateCardRequestBody(cardConfig))
	suite.ctx.Request = req
	err := suite.ctx.ShouldBindJSON(suite.dto)

	verr, _ := err.(validator.ValidationErrors)
	suite.Equal("Title", verr[0].Field())
	suite.Equal("required", verr[0].Tag())
}

func (suite *CardDtoTestSuite) TestBadValidationWithTitleMax100() {
	cardConfig := &factory.CardConfig{Title: strings.Repeat("a", 101)}
	req := httptest.NewRequest("POST", "/", factory.CreateCardRequestBody(cardConfig))
	suite.ctx.Request = req
	err := suite.ctx.ShouldBindJSON(suite.dto)

	verr, _ := err.(validator.ValidationErrors)
	suite.Equal("Title", verr[0].Field())
	suite.Equal("max", verr[0].Tag())
}

func (suite *CardDtoTestSuite) TestBadValidationWithIndexGTE0() {
	req := httptest.NewRequest("POST", "/", factory.CreateCardRequestBody(&factory.CardConfig{Index: -1}))
	suite.ctx.Request = req
	err := suite.ctx.ShouldBindJSON(suite.dto)

	verr, _ := err.(validator.ValidationErrors)
	suite.Equal("Index", verr[0].Field())
	suite.Equal("gte", verr[0].Tag())
}

func (suite *CardDtoTestSuite) TestTransferMethod() {
	dto := factory.NewDtoCard(&factory.CardConfig{})
	var card model.Card
	dto.Transfer(&card)

	suite.Equal(dto.Title, card.Title)
	suite.Equal(dto.Index, card.Index)
}
