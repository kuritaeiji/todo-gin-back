package controller_test

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/controller"
	"github.com/kuritaeiji/todo-gin-back/mock_service"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/stretchr/testify/suite"
)

type ListControllerTestSuite struct {
	suite.Suite
	con             controller.ListController
	ctx             *gin.Context
	rec             *httptest.ResponseRecorder
	listServiceMock *mock_service.MockListService
}

func (suite *ListControllerTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func (suite *ListControllerTestSuite) SetupTest() {
	suite.listServiceMock = mock_service.NewMockListService(gomock.NewController(suite.T()))
	suite.con = controller.TestNewListController(suite.listServiceMock)
	suite.rec = httptest.NewRecorder()
	suite.ctx, _ = gin.CreateTestContext(suite.rec)
}

func TestListControllerSuite(t *testing.T) {
	suite.Run(t, new(ListControllerTestSuite))
}

func (suite *ListControllerTestSuite) TestSuccessIndex() {
	lists := make([]model.List, 0, 2)
	for i := 0; i <= 1; i++ {
		iString := strconv.Itoa(i)
		list := model.List{ID: i, Title: iString}
		card := model.Card{ID: i, Title: iString, ListID: list.ID}
		list.Cards = append(list.Cards, card)
		lists = append(lists, list)
	}
	suite.listServiceMock.EXPECT().Index(suite.ctx).Return(lists, nil)
	suite.con.Index(suite.ctx)

	suite.Equal(200, suite.rec.Code)
	var rLists []model.List
	json.Unmarshal(suite.rec.Body.Bytes(), &rLists)
	suite.Equal(lists[0].ID, rLists[0].ID)
	suite.Equal(lists[0].Title, rLists[0].Title)
	suite.Equal(lists[0].Cards[0].ID, rLists[0].Cards[0].ID)
	suite.Equal(lists[0].Cards[0].Title, rLists[0].Cards[0].Title)
	suite.Equal(lists[0].Cards[0].ListID, rLists[0].Cards[0].ListID)

	suite.Equal(lists[1].ID, rLists[1].ID)
	suite.Equal(lists[1].Title, rLists[1].Title)
	suite.Equal(lists[1].Cards[0].ID, rLists[1].Cards[0].ID)
	suite.Equal(lists[1].Cards[0].Title, rLists[1].Cards[0].Title)
	suite.Equal(lists[1].Cards[0].ListID, rLists[1].Cards[0].ListID)
}

func (suite *ListControllerTestSuite) TestBadIndexWithError() {
	err := errors.New("error")
	suite.listServiceMock.EXPECT().Index(suite.ctx).Return([]model.List{}, err)
	suite.con.Index(suite.ctx)

	suite.Equal(500, suite.rec.Code)
}

func (suite *ListControllerTestSuite) TestSuccessCreate() {
	var list model.List
	list.Title = "test"
	suite.listServiceMock.EXPECT().Create(suite.ctx).Return(list, nil)
	suite.con.Create(suite.ctx)

	suite.Equal(200, suite.rec.Code)
	var rList model.List
	json.Unmarshal(suite.rec.Body.Bytes(), &rList)
	suite.Equal(list.Title, rList.Title)
}

func (suite *ListControllerTestSuite) TestBadCreateWithValidationError() {
	verr := validator.ValidationErrors{}
	suite.listServiceMock.EXPECT().Create(suite.ctx).Return(model.List{}, verr)
	suite.con.Create(suite.ctx)
	suite.Equal(config.ValidationErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.ValidationErrorResponse.Json["content"])
}

func (suite *ListControllerTestSuite) TestBadCreateWithDBError() {
	err := errors.New("DB error")
	suite.listServiceMock.EXPECT().Create(suite.ctx).Return(model.List{}, err)
	suite.con.Create(suite.ctx)

	suite.Equal(500, suite.rec.Code)
}

func (suite *ListControllerTestSuite) TestSuccessUpdate() {
	var list model.List
	suite.listServiceMock.EXPECT().Update(suite.ctx).Return(list, nil)
	suite.con.Update(suite.ctx)

	suite.Equal(200, suite.rec.Code)
	var rList model.List
	json.Unmarshal(suite.rec.Body.Bytes(), &rList)
	suite.Equal(list, rList)
}

func (suite *ListControllerTestSuite) TestBadUpdateWithValidationError() {
	suite.listServiceMock.EXPECT().Update(suite.ctx).Return(model.List{}, validator.ValidationErrors{})
	suite.con.Update(suite.ctx)

	suite.Equal(config.ValidationErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.ValidationErrorResponse.Json["content"])
}

func (suite *ListControllerTestSuite) TestSuccessDestroy() {
	suite.listServiceMock.EXPECT().Destroy(suite.ctx).Return(nil)
	suite.con.Destroy(suite.ctx)

	suite.Equal(200, suite.rec.Code)
}

func (suite *ListControllerTestSuite) TestBadDestroyWithOtherError() {
	suite.listServiceMock.EXPECT().Destroy(suite.ctx).Return(errors.New("error"))
	suite.con.Destroy(suite.ctx)

	suite.Equal(500, suite.rec.Code)
}

func (suite *ListControllerTestSuite) TestSuccessMove() {
	suite.listServiceMock.EXPECT().Move(suite.ctx).Return(nil)
	suite.con.Move(suite.ctx)

	suite.Equal(200, suite.rec.Code)
}

func (suite *ListControllerTestSuite) TestBadMoveOtherError() {
	suite.listServiceMock.EXPECT().Move(suite.ctx).Return(errors.New("other error"))
	suite.con.Move(suite.ctx)

	suite.Equal(500, suite.rec.Code)
}
