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
