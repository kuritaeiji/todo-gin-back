package controller_test

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/controller"
	"github.com/kuritaeiji/todo-gin-back/factory"
	"github.com/kuritaeiji/todo-gin-back/mock_service"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/stretchr/testify/suite"
)

type CardControllerTestSuite struct {
	suite.Suite
	controller      controller.CardController
	cardServiceMock *mock_service.MockCardService
	rec             *httptest.ResponseRecorder
	ctx             *gin.Context
}

func (suite *CardControllerTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func (suite *CardControllerTestSuite) SetupTest() {
	suite.cardServiceMock = mock_service.NewMockCardService(gomock.NewController(suite.T()))
	suite.controller = controller.TestNewCardController(suite.cardServiceMock)
	suite.rec = httptest.NewRecorder()
	suite.ctx, _ = gin.CreateTestContext(suite.rec)
}

func TestCardController(t *testing.T) {
	suite.Run(t, new(CardControllerTestSuite))
}

func (suite *CardControllerTestSuite) TestSuccessCreate() {
	card := factory.NewCard(&factory.CardConfig{})
	suite.cardServiceMock.EXPECT().Create(suite.ctx).Return(card, nil)
	suite.controller.Create(suite.ctx)

	suite.Equal(200, suite.rec.Code)
	var rCard model.Card
	json.Unmarshal(suite.rec.Body.Bytes(), &rCard)
	suite.Equal(card.Title, rCard.Title)
	suite.Equal(card.ID, rCard.ID)
	suite.Equal(card.ListID, rCard.ListID)
}

func (suite *CardControllerTestSuite) TestBadCreateWithValidationError() {
	suite.cardServiceMock.EXPECT().Create(suite.ctx).Return(model.Card{}, validator.ValidationErrors{})
	suite.controller.Create(suite.ctx)

	suite.Equal(config.ValidationErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.ValidationErrorResponse.Json["content"])
}

func (suite *CardControllerTestSuite) TestBadCreateWithOtherError() {
	suite.cardServiceMock.EXPECT().Create(suite.ctx).Return(model.Card{}, errors.New("other error"))
	suite.controller.Create(suite.ctx)

	suite.Equal(500, suite.rec.Code)
}

func (suite *CardControllerTestSuite) TestSuccessUpdate() {
	card := factory.NewCard(&factory.CardConfig{})
	suite.cardServiceMock.EXPECT().Update(suite.ctx).Return(card, nil)
	suite.controller.Update(suite.ctx)

	suite.Equal(200, suite.rec.Code)
	var rCard model.Card
	json.Unmarshal(suite.rec.Body.Bytes(), &rCard)
	suite.Equal(card.Title, rCard.Title)
	suite.Equal(card.ID, rCard.ID)
	suite.Equal(card.ListID, rCard.ListID)
}

func (suite *CardControllerTestSuite) TestBadUpdateWithValidationError() {
	suite.cardServiceMock.EXPECT().Update(suite.ctx).Return(model.Card{}, validator.ValidationErrors{})
	suite.controller.Update(suite.ctx)

	suite.Equal(config.ValidationErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.ValidationErrorResponse.Json["content"])
}

func (suite *CardControllerTestSuite) TestBadUpdateWithOtherError() {
	suite.cardServiceMock.EXPECT().Update(suite.ctx).Return(model.Card{}, errors.New("other error"))
	suite.controller.Update(suite.ctx)

	suite.Equal(500, suite.rec.Code)
}

func (suite *CardControllerTestSuite) TestSuccessDestroyCard() {
	suite.cardServiceMock.EXPECT().Destroy(suite.ctx).Return(nil)
	suite.controller.Destroy(suite.ctx)

	suite.Equal(200, suite.rec.Code)
}

func (suite *CardControllerTestSuite) TestBadDestroyCardWithError() {
	err := errors.New("error")
	suite.cardServiceMock.EXPECT().Destroy(suite.ctx).Return(err)
	suite.controller.Destroy(suite.ctx)

	suite.Equal(500, suite.rec.Code)
}

func (suite *CardControllerTestSuite) TestSuccessMoveCard() {
	suite.cardServiceMock.EXPECT().Move(suite.ctx).Return(nil)
	suite.controller.Move(suite.ctx)

	suite.Equal(200, suite.rec.Code)
}

func (suite *CardControllerTestSuite) TestBadMoveCard() {
	suite.cardServiceMock.EXPECT().Move(suite.ctx).Return(errors.New("db error"))
	suite.controller.Move(suite.ctx)

	suite.Equal(500, suite.rec.Code)
}
