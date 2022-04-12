package repository_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/db"
	"github.com/kuritaeiji/todo-gin-back/factory"
	"github.com/kuritaeiji/todo-gin-back/repository"
	"github.com/stretchr/testify/suite"
)

type ListRepositoryTestSuite struct {
	suite.Suite
	repository repository.ListRepository
}

func (suite *ListRepositoryTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	config.Init()
	db.Init()

	suite.repository = repository.NewListRepository()
}

func (suite *ListRepositoryTestSuite) TearDownTest() {
	db.DeleteAll()
}

func (suite *ListRepositoryTestSuite) TearDownSuite() {
	db.CloseDB()
}

func TestListRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ListRepositoryTestSuite))
}

func (suite *ListRepositoryTestSuite) TestSuccessCreate() {
	user := factory.CreateUser(&factory.UserConfig{})
	list := factory.NewList(&factory.ListConfig{})
	err := suite.repository.Create(&user, &list)

	suite.Nil(err)
	suite.Equal(list.Title, user.Lists[0].Title)
}
