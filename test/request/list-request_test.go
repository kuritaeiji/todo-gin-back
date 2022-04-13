package request_test

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/controller"
	"github.com/kuritaeiji/todo-gin-back/db"
	"github.com/kuritaeiji/todo-gin-back/factory"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/repository"
	"github.com/kuritaeiji/todo-gin-back/server"
	"github.com/kuritaeiji/todo-gin-back/validators"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type ListRequestTestSuite struct {
	suite.Suite
	router     *gin.Engine
	rec        *httptest.ResponseRecorder
	db         *gorm.DB
	repository repository.ListRepository
}

func (suite *ListRequestTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	config.Init()
	validators.Init()
	db.Init()
	suite.router = server.RouterSetup(controller.NewUserController())
	suite.db = db.GetDB()
	suite.repository = repository.NewListRepository()
}

func (suite *ListRequestTestSuite) SetupTest() {
	suite.rec = httptest.NewRecorder()
}

func (suite *ListRequestTestSuite) TearDownTest() {
	db.DeleteAll()
}

func (suite *ListRequestTestSuite) TearDownSuite() {
	db.CloseDB()
}

func TestListRequestSuite(t *testing.T) {
	suite.Run(t, new(ListRequestTestSuite))
}

func (suite *ListRequestTestSuite) TestSuccessIndex() {
	user := factory.CreateUser(&factory.UserConfig{})
	token := factory.CreateAccessToken(user)
	list1 := factory.CreateList(&factory.ListConfig{Index: 0}, user)
	list2 := factory.CreateList(&factory.ListConfig{Index: 1}, user)
	req := httptest.NewRequest("GET", "/api/lists", nil)
	req.Header.Add(config.TokenHeader, token)
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(200, suite.rec.Code)
	var lists []model.List
	json.Unmarshal(suite.rec.Body.Bytes(), &lists)
	suite.Equal([]model.List{list1, list2}, lists)
}

func (suite *ListRequestTestSuite) TestSuccessCreate() {
	user := factory.CreateUser(&factory.UserConfig{})
	token := factory.CreateAccessToken(user)
	listConfig := &factory.ListConfig{}
	body := factory.CreateListRequestBody(listConfig)
	req := httptest.NewRequest("POST", "/api/lists", body)
	req.Header.Add(config.TokenHeader, token)
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(200, suite.rec.Code)
	var list model.List
	json.Unmarshal(suite.rec.Body.Bytes(), &list)
	suite.Equal(listConfig.Title, list.Title)
	suite.Equal(listConfig.Index, list.Index)

	var count int64
	suite.db.Model(&model.List{}).Count(&count)
	suite.Equal(int64(1), count)
}

func (suite *ListRequestTestSuite) TestBadCreateWithValidationError() {
	user := factory.CreateUser(&factory.UserConfig{})
	token := factory.CreateAccessToken(user)
	body := factory.CreateListRequestBody(&factory.ListConfig{NotUseDefaultValue: true})
	req := httptest.NewRequest("POST", "/api/lists", body)
	req.Header.Add(config.TokenHeader, token)
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(config.ValidationErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.ValidationErrorResponse.Json["content"])
}

func (suite *ListRequestTestSuite) TestSuccessUpdate() {
	user := factory.CreateUser(&factory.UserConfig{})
	token := factory.CreateAccessToken(user)
	list := factory.CreateList(&factory.ListConfig{}, user)
	updatingListConfig := &factory.ListConfig{Title: "test title"}
	body := factory.CreateListRequestBody(updatingListConfig)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/lists/%v", list.ID), body)
	req.Header.Add(config.TokenHeader, token)
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(200, suite.rec.Code)
	var rList model.List
	json.Unmarshal(suite.rec.Body.Bytes(), &rList)
	suite.Equal(list.ID, rList.ID)
	suite.Equal(updatingListConfig.Title, rList.Title)
}

func (suite *ListRequestTestSuite) TestBadUpdateWithListRecordNotFound() {
	user := factory.CreateUser(&factory.UserConfig{})
	token := factory.CreateAccessToken(user)
	req := httptest.NewRequest("PUT", "/api/lists/1", nil)
	req.Header.Add(config.TokenHeader, token)
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(config.RecordNotFoundErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.RecordNotFoundErrorResponse.Json["content"])
}

func (suite *ListRequestTestSuite) TestBadUpdateWithForbidden() {
	user := factory.CreateUser(&factory.UserConfig{})
	token := factory.CreateAccessToken(user)
	otherUser := factory.CreateUser(&factory.UserConfig{})
	list := factory.CreateList(&factory.ListConfig{}, otherUser)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/lists/%v", list.ID), nil)
	req.Header.Add(config.TokenHeader, token)
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(config.ForbiddenErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.ForbiddenErrorResponse.Json["content"])
}

func (suite *ListRequestTestSuite) TestBadUpdateWithValidationError() {
	user := factory.CreateUser(&factory.UserConfig{})
	token := factory.CreateAccessToken(user)
	list := factory.CreateList(&factory.ListConfig{}, user)
	body := factory.CreateListRequestBody(&factory.ListConfig{NotUseDefaultValue: true})
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/lists/%v", list.ID), body)
	req.Header.Add(config.TokenHeader, token)
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(config.ValidationErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.ValidationErrorResponse.Json["content"])
}

func (suite *ListRequestTestSuite) TestSuccessDestroy() {
	user := factory.CreateUser(&factory.UserConfig{})
	token := factory.CreateAccessToken(user)
	list := factory.CreateList(&factory.ListConfig{}, user)
	list2 := factory.CreateList(&factory.ListConfig{Index: 1}, user)
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/lists/%v", list.ID), nil)
	req.Header.Add(config.TokenHeader, token)
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(200, suite.rec.Code)
	_, err := suite.repository.Find(list.ID)
	suite.Equal(gorm.ErrRecordNotFound, err)
	rList2, _ := suite.repository.Find(list2.ID)
	suite.Equal(list2.Index-1, rList2.Index)
}

func (suite *ListRequestTestSuite) TestBadDestroyWithNotRecordFound() {
	user := factory.CreateUser(&factory.UserConfig{})
	token := factory.CreateAccessToken(user)
	req := httptest.NewRequest("DELETE", "/api/lists/1", nil)
	req.Header.Add(config.TokenHeader, token)
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(config.RecordNotFoundErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.RecordNotFoundErrorResponse.Json["content"])
}

func (suite *ListRequestTestSuite) TestBadDestroyWithForbiddenError() {
	user := factory.CreateUser(&factory.UserConfig{})
	token := factory.CreateAccessToken(user)
	otherUser := factory.CreateUser(&factory.UserConfig{})
	list := factory.CreateList(&factory.ListConfig{}, otherUser)
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/lists/%v", list.ID), nil)
	req.Header.Add(config.TokenHeader, token)
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(config.ForbiddenErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.ForbiddenErrorResponse.Json["content"])
}
