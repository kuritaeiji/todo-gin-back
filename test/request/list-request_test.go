package request_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/controller"
	"github.com/kuritaeiji/todo-gin-back/db"
	"github.com/kuritaeiji/todo-gin-back/factory"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/server"
	"github.com/kuritaeiji/todo-gin-back/validators"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type ListRequestTestSuite struct {
	suite.Suite
	router *gin.Engine
	rec    *httptest.ResponseRecorder
	db     *gorm.DB
}

func (suite *ListRequestTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	config.Init()
	validators.Init()
	db.Init()
	suite.router = server.RouterSetup(controller.NewUserController())
	suite.db = db.GetDB()
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
